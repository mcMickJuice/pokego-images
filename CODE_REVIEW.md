# Code Review: pokego-images (Final)

**Reviewer**: Senior Engineer
**Date**: January 12, 2026
**Branch**: `mj/separate-executables`
**Revision**: 3

---

## Summary of Changes Since Last Review

You addressed all high and medium priority issues:

- ✅ Fixed resource leak: `defer resp.Body.Close()` now immediately follows error check
- ✅ CLI now handles the error from `Write()`
- ✅ Fixed typo: `pokemonSprintUrl` → `pokemonSpriteURL`
- ✅ Webserver logging now uses `log.Fatalf`
- ✅ Defer pattern is now consistent across both HTTP functions
- ✅ Fixed comment typo: "soas" → "so as"
- ✅ Server address is now configurable via parameter
- ✅ Webserver refactored to struct with constructor pattern (nice improvement!)

The codebase is in good shape. Below are the remaining minor issues and a couple of new observations.

---

## Remaining Issues

### Issue 1: Square image assumption (`pokeimage/image.go:58-63`)

```go
maxBounds := pi.Bounds().Max.X
for r := 0; r < maxBounds; r++ {
    for c := 0; c < maxBounds; c++ {
```

**Why this matters**: Uses `Max.X` for both width and height. Works for square Pokemon sprites but would produce incorrect output for rectangular images.

**Safer pattern**:
```go
bounds := pi.Bounds()
for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
        grayscale := toGrayscale(pi.At(x, y))
```

**Priority**: Low - Pokemon sprites are square, so this works for the current use case.

---

### Issue 2: No tests

Still no `*_test.go` files. Consider adding tests for the pure functions in `pokeimage` as a starting point.

**Priority**: Low for a personal project, but recommended before sharing or deploying.

---

### Issue 3: Naming inconsistency (`webserver/handler.go:11,15`)

```go
type PokemonWebServer struct {  // "WebServer" with capital S
    addr string
}

func NewPokemonWebserver(addr string) PokemonWebServer {  // "Webserver" with lowercase s
```

**Why this matters**: Go convention is to keep naming consistent. Since `Server` is a complete word (not an initialism), either `WebServer` or `Webserver` is acceptable, but pick one.

**Fix**: Rename to match - either both `PokemonWebServer`/`NewPokemonWebServer` or both `PokemonWebserver`/`NewPokemonWebserver`.

**Priority**: Low - cosmetic.

---

### Issue 4: Mixed logging in webserver (`webserver/handler.go:42,46`)

```go
fmt.Printf("webserver started at %s\n", s.addr)  // line 46
// ...
fmt.Printf("error writing to response: %v", err)  // line 42
```

You updated `cmd/webserver/main.go` to use `log`, but the handler still uses `fmt.Printf`. For consistency, consider using `log.Printf` throughout.

**Why this matters**: `log.Printf` includes timestamps and writes to stderr, which is more appropriate for operational messages and errors.

**Priority**: Low - cosmetic.

---

## Code Quality Summary

| Category | Status |
|----------|--------|
| Critical bugs | None |
| Error handling | Good |
| Resource management | Good |
| Naming conventions | Good (minor inconsistency) |
| Code organization | Good |
| Logging | Acceptable (minor inconsistency) |
| Tests | Missing |

---

## What You Did Well

1. **Project structure** - Clean separation with `cmd/`, `internal/`, and logical package boundaries
2. **Error handling** - Errors are properly checked, wrapped with context, and propagated
3. **Sentinel errors** - `ErrPokemonNotFound` allows callers to handle specific error cases with `errors.Is()`
4. **HTTP status codes** - Webserver returns appropriate 404 vs 500 responses
5. **Interface usage** - `Write(w io.Writer)` accepts any writer, making the code flexible and testable
6. **Struct pattern for webserver** - The refactor to `PokemonWebServer` with `NewPokemonWebserver()` and `Start()` is idiomatic Go

---

## Suggested Next Steps

1. **Add basic tests** - Even a few tests for `blankLine()` and `grayscaleToAscii()` would catch regressions
2. **Consider graceful shutdown** - The webserver could handle `SIGTERM`/`SIGINT` for clean container deployments
3. **Add request logging** - Log each request for observability
4. **Consider structured logging** - For production, packages like `slog` (stdlib in Go 1.21+) or `zerolog` provide structured JSON logs

---

## Final Assessment

This is solid, idiomatic Go code. You've addressed all the significant issues from earlier reviews. The remaining items are cosmetic or nice-to-haves. The code is ready to merge.

Good work learning Go - you've picked up the idioms quickly.
