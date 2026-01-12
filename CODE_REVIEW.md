# Code Review: pokego-images (Updated)

**Reviewer**: Senior Engineer
**Date**: January 12, 2026
**Branch**: `mj/separate-executables`
**Revision**: 2 (post-fixes)

---

## Summary of Changes Since Last Review

You addressed several key issues:

- ✅ Replaced `panic` with `log.Fatalf` in CLI
- ✅ Fixed critical bug: error check now comes before using `resp` in `getPokemon()`
- ✅ Constants renamed to Go conventions (`rFactor`, `gFactor`, etc.)
- ✅ Magic number `32` replaced with named constant `asciiSpaceCharCode`
- ✅ Removed unused `ALL_POKEMON_URL` constant
- ✅ `PokemonImage.Write()` now returns an error
- ✅ Removed unnecessary blank line in `NewPokemonClient`
- ✅ Webserver handler now checks error from `Write()`

Good work on these fixes. Below are the remaining issues.

---

## Remaining Issues

### Issue 1: Resource leak - defer placement (`client.go:63-71`)

```go
func getPokemon(pokemonName string) (pokemonResponse, error) {
    resp, err := http.Get(fmt.Sprintf(pokemonDetailURL, pokemonName))

    if err != nil {
        return pokemonResponse{}, fmt.Errorf("unknown error fetching pokemon: %w", err)
    }

    if resp.StatusCode == http.StatusNotFound {
        return pokemonResponse{}, ErrPokemonNotFound  // Body not closed!
    }

    if resp.StatusCode != http.StatusOK {
        return pokemonResponse{}, fmt.Errorf(...)  // Body not closed!
    }

    defer resp.Body.Close()  // Too late - early returns above leak the body
```

**Why this matters**: If the API returns a 404 or any non-200 status, you return early *without* closing the response body. This leaks file descriptors and can cause resource exhaustion under load.

**Fix**: Move the `defer` immediately after the error check:
```go
if err != nil {
    return pokemonResponse{}, fmt.Errorf("unknown error fetching pokemon: %w", err)
}
defer resp.Body.Close()  // Right here, before any other returns

if resp.StatusCode == http.StatusNotFound {
    return pokemonResponse{}, ErrPokemonNotFound
}
```

---

### Issue 2: CLI ignores Write error (`cmd/cli/main.go:21`)

```go
pokeimage.NewPokemonImage(image).Write(os.Stdout)
```

**Why this matters**: You updated `Write()` to return an error (good!), but the CLI doesn't check it. If writing to stdout fails (broken pipe, disk full, etc.), the error is silently ignored.

**Fix**:
```go
if err := pokeimage.NewPokemonImage(image).Write(os.Stdout); err != nil {
    log.Fatalf("Failed to write image: %v", err)
}
```

---

### Issue 3: Typo in parameter name (`client.go:82`)

```go
func getPokemonSprite(pokemonSprintUrl string) (image.Image, error) {
                            ^^^^^
```

`Sprint` should be `Sprite`, and `Url` should be `URL` per Go conventions for initialisms.

**Fix**: `pokemonSpriteURL`

---

### Issue 4: Webserver main.go logging issues (`cmd/webserver/main.go:9-14`)

```go
fmt.Println("starting webserver")
err := webserver.CreateWebserver()
if err != nil {
    fmt.Printf("error starting webserver: %v", err)  // No newline, goes to stdout
}
fmt.Println("exiting webserver")  // Prints even on error
```

**Issues**:
1. Error message missing newline
2. Error goes to stdout instead of stderr
3. "exiting webserver" prints even when there's an error (misleading)

**Fix**:
```go
func main() {
    log.Println("starting webserver on localhost:8080")
    if err := webserver.CreateWebserver(); err != nil {
        log.Fatalf("webserver error: %v", err)
    }
}
```

Note: `log.Fatalf` prints to stderr and exits with code 1, so nothing after it runs.

---

### Issue 5: Inconsistent defer patterns (`client.go:71` vs `client.go:89-94`)

In `getPokemon`:
```go
defer resp.Body.Close()
```

In `getPokemonSprite`:
```go
defer func() {
    err := resp.Body.Close()
    if err != nil {
        fmt.Printf("failed to close response body: %v", err)
    }
}()
```

**Why this matters**: Inconsistency makes code harder to maintain. Pick one pattern and stick with it.

**Recommendation**: Use the simple `defer resp.Body.Close()` in both places. For HTTP response bodies, the Close error is almost never actionable - the data has already been read. The verbose pattern adds noise without benefit.

---

### Issue 6: Comment typo (`pokeimage/image.go:73`)

```go
// append newline here instead of above soas not to break blaneLine logic
                                      ^^^^
```

"soas" should be "so as".

---

### Issue 7: Hardcoded server address (`webserver/handler.go:37`)

```go
err := http.ListenAndServe("localhost:8080", mux)
```

**Why this matters**: Hardcoded configuration makes the code inflexible. You can't run on a different port without changing source code.

**Fix**: Accept the address as a parameter:
```go
func CreateWebserver(addr string) error {
    // ...
    return http.ListenAndServe(addr, mux)
}
```

Then in main:
```go
webserver.CreateWebserver(":8080")
```

Or read from environment variable / flag for production flexibility.

---

### Issue 8: Square image assumption (`pokeimage/image.go:58-63`)

```go
maxBounds := pi.Bounds().Max.X
for r := 0; r < maxBounds; r++ {
    for c := 0; c < maxBounds; c++ {
```

**Why this matters**: Uses `Max.X` for both width and height. Works for square Pokemon sprites but breaks for rectangular images.

**Safer pattern**:
```go
bounds := pi.Bounds()
for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
        grayscale := toGrayscale(pi.At(x, y))
```

Note: Using `Min.X`/`Min.Y` as start points is important because some image formats have non-zero origins.

---

### Issue 9: No tests

Still no `*_test.go` files in the project.

**Recommendation**: Start with testing pure functions that are easy to verify:

```go
// internal/pokeimage/image_test.go
package pokeimage

import (
    "bytes"
    "image"
    "image/color"
    "testing"
)

func TestBlankLine(t *testing.T) {
    tests := []struct {
        name  string
        input []byte
        want  bool
    }{
        {"all spaces", []byte("        "), true},
        {"has content", []byte("  @  "), false},
        {"empty", []byte{}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := blankLine(tt.input); got != tt.want {
                t.Errorf("blankLine(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}

func TestGrayscaleToAscii(t *testing.T) {
    // Test boundary conditions
    if got := grayscaleToAscii(0); got != ' ' {
        t.Errorf("grayscaleToAscii(0) = %q, want ' '", got)
    }
}
```

Run with: `go test ./...`

---

## Summary of Remaining Action Items

### High Priority
1. **Fix resource leak**: Move `defer resp.Body.Close()` before status code checks
2. **Handle Write error in CLI**: Check and handle the error from `Write()`

### Medium Priority
3. Fix typo: `pokemonSprintUrl` → `pokemonSpriteURL`
4. Fix webserver logging (use `log.Fatalf`, remove "exiting" message)
5. Make defer pattern consistent across both HTTP functions

### Low Priority
6. Fix comment typo "soas" → "so as"
7. Make server address configurable
8. Handle non-square images
9. Add tests

---

## Progress

| Category | Original Issues | Fixed | Remaining |
|----------|----------------|-------|-----------|
| Critical | 1 | 1 | 0 |
| High Priority | 3 | 2 | 2 |
| Medium Priority | 4 | 2 | 3 |
| Low Priority | 5 | 1 | 4 |

You're making good progress. The critical nil-pointer bug is fixed, which was the most important item. The remaining high-priority issues (resource leak and unhandled error) are worth addressing before merging.
