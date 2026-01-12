# Code Review: pokego-images

**Reviewer**: Senior Engineer
**Date**: January 12, 2026
**Branch**: `mj/separate-executables`

---

## Executive Summary

This is a well-organized Go project that fetches Pokemon sprites from the PokeAPI and converts them to ASCII art. The code is generally clean and demonstrates good instincts around project structure. There are several areas where adopting more idiomatic Go patterns would improve maintainability, testability, and robustness.

**Overall Assessment**: Solid foundation with room for improvement in error handling, interface design, and Go conventions.

---

## Project Structure

### What's Good

Your project layout follows the widely-adopted [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
.
├── cmd/
│   ├── cli/main.go
│   └── webserver/main.go
├── internal/
│   ├── pokeimage/
│   ├── pokemon/
│   └── webserver/
├── go.mod
└── .gitignore
```

- **`cmd/`** for executables is correct
- **`internal/`** prevents external packages from importing your code, which is appropriate here
- Separation of concerns between `pokemon` (API client), `pokeimage` (image processing), and `webserver` is logical

### Suggestions

1. **Consider a `pkg/` directory** if you ever want to expose packages for external use. Right now everything is in `internal/`, which is fine for a self-contained application.

2. **Add a `Makefile` or `justfile`** for common operations:
   ```makefile
   build:
       go build -o bin/cli ./cmd/cli
       go build -o bin/webserver ./cmd/webserver
   ```

---

## Error Handling

This is the area with the most room for improvement. Go has specific idioms around error handling that differ from exception-based languages.

### Issue 1: Checking error *after* using response (`client.go:56-69`)

```go
func getPokemon(pokemonName string) (pokemonResponse, error) {
    resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, pokemonName))

    if resp.StatusCode == http.StatusNotFound {  // Using resp before checking err!
        return pokemonResponse{}, ErrPokemonNotFound
    }
    // ...
    if err != nil {  // Too late - we already used resp above
        return pokemonResponse{}, fmt.Errorf("unknown error fetching pokemon: %w", err)
    }
```

**Why this matters**: If `http.Get` returns an error, `resp` may be `nil`, causing a panic when you access `resp.StatusCode`. The error check must come first.

**Idiomatic pattern**:
```go
resp, err := http.Get(url)
if err != nil {
    return pokemonResponse{}, fmt.Errorf("failed to fetch pokemon: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode == http.StatusNotFound {
    return pokemonResponse{}, ErrPokemonNotFound
}
```

### Issue 2: Using `panic` in CLI (`cmd/cli/main.go:17`)

```go
if err != nil {
    panic(err)
}
```

**Why this matters**: `panic` is reserved for truly unrecoverable situations (programmer errors, invariant violations). User input errors are expected and should be handled gracefully.

**Idiomatic pattern**:
```go
if err != nil {
    fmt.Fprintf(os.Stderr, "error: %v\n", err)
    os.Exit(1)
}
```

Or use `log.Fatal(err)` which does both.

### Issue 3: Silent error handling (`pokeimage/image.go:74-78`)

```go
_, err := w.Write(line)
if err != nil {
    fmt.Printf("error writing to writer: %v", err)
    return  // Silently returns, caller doesn't know there was an error
}
```

**Why this matters**: The `Write` method signature is `func (pi PokemonImage) Write(w io.Writer)` - it returns nothing. The caller has no way to know if writing failed. This violates the Go principle of explicit error handling.

**Idiomatic pattern**: Return the error and let the caller decide:
```go
func (pi PokemonImage) Write(w io.Writer) error {
    // ...
    if _, err := w.Write(line); err != nil {
        return fmt.Errorf("failed to write line: %w", err)
    }
    // ...
    return nil
}
```

### Issue 4: Inconsistent defer pattern for `resp.Body.Close()`

You have this pattern repeated:
```go
defer func() {
    err := resp.Body.Close()
    if err != nil {
        fmt.Printf("failed to close response body: %v", err)
    }
}()
```

**Why this matters**: While checking `Close()` errors is thorough, for HTTP response bodies it's typically overkill since you've already read the data. The standard pattern is simply:
```go
defer resp.Body.Close()
```

If you do want to handle the error (which can matter for writable file handles), use `errcheck` linter or a helper function rather than inline anonymous functions.

---

## Code Style & Idiomatic Go

### Issue 5: Constant naming (`pokeimage/image.go:22-24`, `client.go:11-12`)

```go
const R_FACTOR float32 = 0.299
const ALL_POKEMON_URL = "https://pokeapi.co/api/v2/pokemon?limit=400"
```

**Why this matters**: Go uses `MixedCaps` or `mixedCaps`, not `SCREAMING_SNAKE_CASE`. The latter is from C/Java traditions. In Go, the case of the first letter determines visibility (exported vs unexported).

**Idiomatic pattern**:
```go
const (
    rFactor float32 = 0.299  // unexported
    gFactor float32 = 0.587
    bFactor float32 = 0.114
)

const (
    allPokemonURL   = "https://pokeapi.co/api/v2/pokemon?limit=400"
    pokemonDetailURL = "https://pokeapi.co/api/v2/pokemon/%s"
)
```

Note: `URL` not `Url` per Go conventions for initialisms.

### Issue 6: Unnecessary blank lines in functions

```go
func NewPokemonClient(pokemonName string) *PokemonClient {

    return &PokemonClient{
```

**Why this matters**: Go style (enforced by `gofmt`) doesn't use blank lines at the start of function bodies. Run `gofmt -w .` to auto-fix these.

### Issue 7: Magic numbers (`pokeimage/image.go:35`)

```go
if char != 32 {
```

**Why this matters**: What is 32? (It's ASCII space.) Magic numbers make code harder to understand.

**Idiomatic pattern**:
```go
if char != ' ' {  // Can compare byte to rune literal
```

### Issue 8: Receiver naming (`client.go:39`)

```go
func (pc PokemonClient) GetPokemonSprite() (image.Image, error) {
```

**Why this matters**: Go convention is single-letter or very short receiver names. Multi-letter abbreviations are fine for clarity but should be consistent.

**Idiomatic pattern**: Either `p` for pokemon or `c` for client, used consistently:
```go
func (c PokemonClient) GetPokemonSprite() (image.Image, error) {
```

### Issue 9: Typo in function name (`client.go:87`)

```go
func getPokemonSprite(pokemonSprintUrl string) (image.Image, error) {
                            ^^^^^^
```

`Sprint` should be `Sprite`. Also `Url` should be `URL` per Go conventions.

---

## Design & Architecture

### Issue 10: `PokemonClient` design

The current design creates a client per Pokemon:
```go
pokemonClient := pokemon.NewPokemonClient("snorlax")
image, err := pokemonClient.GetPokemonSprite()
```

**Why this matters**: This is unusual. A "client" typically represents a connection or configuration that can be reused. Creating a new client per request adds conceptual overhead.

**Alternative 1** - Stateless function:
```go
image, err := pokemon.GetSprite("snorlax")
```

**Alternative 2** - Reusable client with method parameter:
```go
client := pokemon.NewClient()  // Could configure timeout, base URL, etc.
image, err := client.GetSprite("snorlax")
```

Alternative 2 is better for testability (you can inject a mock client) and extensibility (add caching, rate limiting, etc.).

### Issue 11: Hardcoded server address (`webserver/handler.go:36`)

```go
err := http.ListenAndServe("localhost:8080", mux)
```

**Why this matters**: Configuration should be injectable for flexibility and testing.

**Idiomatic pattern**:
```go
func CreateWebserver(addr string) error {
    // ...
    return http.ListenAndServe(addr, mux)
}
```

Or accept a config struct for multiple settings.

### Issue 12: Square image assumption (`pokeimage/image.go:57-62`)

```go
maxBounds := pi.Bounds().Max.X
for r := 0; r < maxBounds; r++ {
    for c := 0; c < maxBounds; c++ {
```

**Why this matters**: You use `Max.X` for both dimensions, assuming the image is square. This works for Pokemon sprites but would break for rectangular images.

**Safer pattern**:
```go
bounds := pi.Bounds()
for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
```

---

## Testing

### Issue 13: No tests

**Why this matters**: Go has excellent built-in testing. Tests are expected in Go projects and live alongside the code as `*_test.go` files.

**Recommendation**: Add at minimum:
- `internal/pokemon/client_test.go` - Test error handling paths
- `internal/pokeimage/image_test.go` - Test ASCII conversion with known inputs

Example test structure:
```go
// image_test.go
package pokeimage

import "testing"

func TestGrayscaleToAscii(t *testing.T) {
    tests := []struct {
        brightness float32
        want       byte
    }{
        {0, ' '},
        {65535, '@'},
    }
    for _, tt := range tests {
        got := grayscaleToAscii(tt.brightness)
        if got != tt.want {
            t.Errorf("grayscaleToAscii(%v) = %v, want %v", tt.brightness, got, tt.want)
        }
    }
}
```

---

## Minor Issues

### Issue 14: Unused constant (`client.go:11`)

```go
const ALL_POKEMON_URL = "https://pokeapi.co/api/v2/pokemon?limit=400"
```

This constant is defined but never used. Remove dead code.

### Issue 15: Webserver logging (`cmd/webserver/main.go:9-14`)

```go
fmt.Println("starting webserver")
err := webserver.CreateWebserver()
if err != nil {
    fmt.Printf("error starting webserver: %v", err)
}
fmt.Println("exiting webserver")
```

**Issues**:
- `fmt.Printf` for errors should go to `os.Stderr`
- "exiting webserver" will print even on error
- Missing newline in error message

**Better**:
```go
log.Println("starting webserver on :8080")
if err := webserver.CreateWebserver(":8080"); err != nil {
    log.Fatalf("webserver error: %v", err)
}
```

### Issue 16: Comment typo (`pokeimage/image.go:72`)

```go
// append newline here instead of above soas not to break blaneLine logic
                                      ^^^^
```

"soas" should be "so as".

---

## Security Considerations

### Issue 17: No input validation on Pokemon name

The Pokemon name from user input (CLI flag or URL path) is passed directly to the API URL:
```go
resp, err := http.Get(fmt.Sprintf(POKEMON_DETAIL_URL, pokemonName))
```

While PokeAPI is trusted and Go's HTTP client handles URL encoding, it's good practice to validate/sanitize input. A malicious actor could potentially craft inputs that cause unexpected behavior.

**Recommendation**: At minimum, validate the Pokemon name matches expected patterns (alphanumeric, hyphens for names like "mr-mime").

---

## Summary of Action Items

### Critical (Fix Now)
1. Check `err` before using `resp` in `getPokemon()`

### High Priority
2. Return errors from `PokemonImage.Write()` instead of printing
3. Replace `panic(err)` with graceful error handling in CLI
4. Add basic tests

### Medium Priority
5. Fix constant naming to Go conventions
6. Fix typo `SprintUrl` -> `SpriteURL`
7. Remove unused `ALL_POKEMON_URL` constant
8. Make server address configurable

### Low Priority (Nice to Have)
9. Add Makefile for builds
10. Use proper logging instead of fmt.Print for webserver
11. Consider redesigning PokemonClient to be reusable
12. Handle non-square images in ASCII conversion

---

## Resources for Learning Idiomatic Go

- [Effective Go](https://go.dev/doc/effective_go) - Official style guide
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Common review feedback
- [Go Proverbs](https://go-proverbs.github.io/) - Philosophy of Go
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) - Practical conventions

---

Good work getting this far. The code does what it's supposed to do, and you've made smart choices about project organization. The issues above are the kind of things that come with experience in the language. Keep building.
