# Go Backend Toolkit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a reusable Go backend toolkit with common utility packages, Chinese GoDoc comments, unit tests, and Chinese usage documentation.

**Architecture:** The project uses small top-level packages grouped by responsibility instead of a single `utils` package. Public APIs use standard-library types where possible, with only stable dependencies for token-bucket rate limiting and YAML parsing. Each package is independently testable and documented.

**Tech Stack:** Go 1.22+, standard library, `golang.org/x/time/rate`, `gopkg.in/yaml.v3`, Apache-2.0 license.

## Global Constraints

- Module path: `github.com/z-nuo/GLTools`.
- Go version: Go 1.22 or higher.
- All exported functions, types, constants, and variables must have Chinese GoDoc comments.
- Do not bind HTTP helpers to Gin, Echo, Fiber, or any other web framework.
- Do not add MySQL, Redis, Kafka, JWT, object storage, dependency injection, or CLI features in this version.
- Public APIs should prefer `context.Context`, `error`, `time.Duration`, `http.Client`, `http.Request`, and `http.Response` where applicable.
- Final verification commands: `go test ./...` and `go vet ./...`.

---

## File Structure

- Create `go.mod`: module declaration and dependencies.
- Modify `README.md`: Chinese project overview, install command, quick start, package list, test commands.
- Create `docs/usage.md`: full Chinese usage guide with package-by-package examples.
- Create `glstrings/strings.go` and `glstrings/strings_test.go`: string helpers.
- Create `gltime/time.go` and `gltime/time_test.go`: time helpers.
- Create `glconv/conv.go` and `glconv/conv_test.go`: conversion helpers.
- Create `glslice/slice.go` and `glslice/slice_test.go`: generic and typed slice helpers.
- Create `glmap/map.go` and `glmap/map_test.go`: generic map helpers.
- Create `glrand/rand.go` and `glrand/rand_test.go`: secure random helpers.
- Create `glcrypto/crypto.go` and `glcrypto/crypto_test.go`: digest, HMAC, Base64 helpers.
- Create `gljson/json.go` and `gljson/json_test.go`: JSON helpers.
- Create `glfile/file.go` and `glfile/file_test.go`: file and path helpers.
- Create `glerror/error.go` and `glerror/error_test.go`: code error helpers.
- Create `glhttp/response.go`, `glhttp/client.go`, and `glhttp/http_test.go`: response model and HTTP client.
- Create `gllog/log.go` and `gllog/log_test.go`: `slog` setup helpers.
- Create `glconfig/config.go` and `glconfig/config_test.go`: JSON/YAML/env configuration helpers.
- Create `glretry/retry.go` and `glretry/retry_test.go`: retry helpers.
- Create `gllimit/limit.go` and `gllimit/limit_test.go`: rate limiter wrapper.
- Create `glid/id.go` and `glid/id_test.go`: Snowflake-style ID generator.

---

### Task 1: Initialize Go Module and Project Baseline

**Files:**
- Create: `go.mod`
- Modify: `README.md`

**Interfaces:**
- Consumes: repository root.
- Produces: module path `github.com/z-nuo/GLTools`, Go 1.22 baseline, dependencies for later packages.

- [ ] **Step 1: Create module file**

```go
module github.com/z-nuo/GLTools

go 1.22

require (
	golang.org/x/time v0.5.0
	gopkg.in/yaml.v3 v3.0.1
)
```

- [ ] **Step 2: Run module tidy**

Run: `go mod tidy`

Expected: command succeeds and creates `go.sum` if dependency checksums are needed.

- [ ] **Step 3: Replace README with Chinese baseline**

````markdown
# GLTools

GLTools 是一个面向 Go 后端服务的通用工具库，封装字符串、时间、类型转换、切片、Map、随机数、加密摘要、JSON、文件、HTTP、日志、配置、错误码、重试、限流和 ID 生成等常用能力。

## 安装

```bash
go get github.com/z-nuo/GLTools
```

## 快速开始

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glstrings"
)

func main() {
	fmt.Println(glstrings.IsBlank("   "))
}
```

## 测试

```bash
go test ./...
go vet ./...
```

## 使用说明

完整说明见 `docs/usage.md`。
````

- [ ] **Step 4: Verify baseline**

Run: `go test ./...`

Expected: PASS or message indicating no packages before package files are added.

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum README.md
git commit -m "chore: initialize go module"
```

---

### Task 2: Basic Utility Packages

**Files:**
- Create: `glstrings/strings.go`
- Create: `glstrings/strings_test.go`
- Create: `gltime/time.go`
- Create: `gltime/time_test.go`
- Create: `glconv/conv.go`
- Create: `glconv/conv_test.go`
- Create: `glslice/slice.go`
- Create: `glslice/slice_test.go`
- Create: `glmap/map.go`
- Create: `glmap/map_test.go`

**Interfaces:**
- Consumes: module path from Task 1.
- Produces:
  - `glstrings.IsBlank(s string) bool`
  - `glstrings.Trim(s string) string`
  - `glstrings.Truncate(s string, maxRunes int) string`
  - `glstrings.SnakeToCamel(s string) string`
  - `glstrings.CamelToSnake(s string) string`
  - `gltime.Format(t time.Time, layout string) string`
  - `gltime.Parse(layout string, value string) (time.Time, error)`
  - `gltime.UnixMilli(t time.Time) int64`
  - `gltime.FromUnixMilli(ms int64) time.Time`
  - `gltime.DayStart(t time.Time) time.Time`
  - `gltime.DayEnd(t time.Time) time.Time`
  - `glconv.ToInt(s string) (int, error)`
  - `glconv.ToIntDefault(s string, def int) int`
  - `glconv.ToInt64(s string) (int64, error)`
  - `glconv.ToFloat64(s string) (float64, error)`
  - `glconv.ToBool(s string) (bool, error)`
  - `glconv.String(v any) string`
  - `glslice.Contains[T comparable](items []T, target T) bool`
  - `glslice.Unique[T comparable](items []T) []T`
  - `glslice.Filter[T any](items []T, keep func(T) bool) []T`
  - `glslice.CompactStrings(items []string) []string`
  - `glmap.Keys[K comparable, V any](m map[K]V) []K`
  - `glmap.Values[K comparable, V any](m map[K]V) []V`
  - `glmap.HasKey[K comparable, V any](m map[K]V, key K) bool`
  - `glmap.Merge[K comparable, V any](left map[K]V, right map[K]V) map[K]V`

- [ ] **Step 1: Write failing tests for `glstrings`**

```go
func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{name: "empty", in: "", want: true},
		{name: "spaces", in: " \t\n", want: true},
		{name: "text", in: "go", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.in); got != tt.want {
				t.Fatalf("IsBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateUsesRunes(t *testing.T) {
	if got := Truncate("你好Go", 2); got != "你好" {
		t.Fatalf("Truncate() = %q, want %q", got, "你好")
	}
}
```

- [ ] **Step 2: Run `glstrings` tests to verify RED**

Run: `go test ./glstrings`

Expected: FAIL because package files or functions are missing.

- [ ] **Step 3: Implement `glstrings` with Chinese GoDoc**

Implement the exported functions listed in this task's interface. Use `strings.TrimSpace`, rune slicing for truncation, and deterministic naming conversion.

- [ ] **Step 4: Run `glstrings` tests to verify GREEN**

Run: `go test ./glstrings`

Expected: PASS.

- [ ] **Step 5: Repeat RED-GREEN for `gltime`**

Write tests covering `Format`, `Parse`, `UnixMilli`, `FromUnixMilli`, `DayStart`, and `DayEnd`. Use a fixed time value:

```go
base := time.Date(2026, 6, 30, 15, 4, 5, 123000000, time.Local)
```

Run: `go test ./gltime`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 6: Repeat RED-GREEN for `glconv`**

Write table tests for valid values, invalid values, and default fallback behavior:

```go
got := ToIntDefault("bad", 7)
if got != 7 {
	t.Fatalf("ToIntDefault() = %d, want 7", got)
}
```

Run: `go test ./glconv`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 7: Repeat RED-GREEN for `glslice`**

Write tests for `Contains`, `Unique`, `Filter`, and `CompactStrings`. Preserve first-seen order in `Unique`.

Run: `go test ./glslice`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 8: Repeat RED-GREEN for `glmap`**

Write tests for keys, values, key existence, and right-side overwrite in `Merge`.

Run: `go test ./glmap`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 9: Run all current tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 10: Commit**

```bash
git add glstrings gltime glconv glslice glmap
git commit -m "feat: add basic utility packages"
```

---

### Task 3: Encoding, Random, JSON, and File Packages

**Files:**
- Create: `glrand/rand.go`
- Create: `glrand/rand_test.go`
- Create: `glcrypto/crypto.go`
- Create: `glcrypto/crypto_test.go`
- Create: `gljson/json.go`
- Create: `gljson/json_test.go`
- Create: `glfile/file.go`
- Create: `glfile/file_test.go`

**Interfaces:**
- Consumes: module baseline from Task 1.
- Produces:
  - `glrand.NumericCode(length int) (string, error)`
  - `glrand.String(length int) (string, error)`
  - `glrand.StringFromCharset(length int, charset string) (string, error)`
  - `glcrypto.MD5Hex(s string) string`
  - `glcrypto.SHA256Hex(s string) string`
  - `glcrypto.HMACSHA256Hex(message string, secret string) string`
  - `glcrypto.Base64Encode(s string) string`
  - `glcrypto.Base64Decode(s string) (string, error)`
  - `gljson.Marshal(v any) ([]byte, error)`
  - `gljson.Unmarshal(data []byte, v any) error`
  - `gljson.Valid(s string) bool`
  - `gljson.Pretty(v any) (string, error)`
  - `glfile.Exists(path string) bool`
  - `glfile.IsFile(path string) bool`
  - `glfile.IsDir(path string) bool`
  - `glfile.EnsureDir(path string) error`
  - `glfile.ReadText(path string) (string, error)`
  - `glfile.WriteText(path string, content string) error`
  - `glfile.Ext(path string) string`
  - `glfile.Join(elem ...string) string`

- [ ] **Step 1: Write failing tests for `glrand`**

Test length, numeric-only output, custom charset, and invalid length:

```go
got, err := NumericCode(6)
if err != nil {
	t.Fatal(err)
}
if len(got) != 6 {
	t.Fatalf("len = %d, want 6", len(got))
}
for _, r := range got {
	if r < '0' || r > '9' {
		t.Fatalf("non numeric rune %q", r)
	}
}
```

Run: `go test ./glrand`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 2: Implement `glrand`**

Use `crypto/rand.Int` with `math/big`. Return errors for `length <= 0` and empty charset. Add Chinese GoDoc for all exported functions.

- [ ] **Step 3: Write failing tests for `glcrypto`**

Use known digests:

```go
if got := MD5Hex("abc"); got != "900150983cd24fb0d6963f7d28e17f72" {
	t.Fatalf("MD5Hex() = %s", got)
}
if got := SHA256Hex("abc"); got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
	t.Fatalf("SHA256Hex() = %s", got)
}
```

Run: `go test ./glcrypto`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 4: Implement `glcrypto`**

Use `crypto/md5`, `crypto/sha256`, `crypto/hmac`, `encoding/hex`, and `encoding/base64`. Add Chinese GoDoc and keep returned digest lowercase hex.

- [ ] **Step 5: Write failing tests for `gljson`**

Test marshal, unmarshal, valid JSON, invalid JSON, and pretty output containing newlines.

Run: `go test ./gljson`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 6: Implement `gljson`**

Use `github.com/json-iterator/go` with `ConfigCompatibleWithStandardLibrary`. `Valid` accepts a string and calls the compatible API's `Valid([]byte(s))`. `Pretty` uses `MarshalIndent(v, "", "  ")`.

- [ ] **Step 7: Write failing tests for `glfile`**

Use `t.TempDir()` and verify directory creation, text write/read, extension lookup, and path joining.

Run: `go test ./glfile`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 8: Implement `glfile`**

Use `os.Stat`, `os.MkdirAll`, `os.ReadFile`, `os.WriteFile`, `path/filepath.Ext`, and `filepath.Join`. `WriteText` creates parent directories when needed.

- [ ] **Step 9: Run all current tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 10: Commit**

```bash
git add glrand glcrypto gljson glfile
git commit -m "feat: add encoding random json and file helpers"
```

---

### Task 4: Error Code and HTTP Packages

**Files:**
- Create: `glerror/error.go`
- Create: `glerror/error_test.go`
- Create: `glhttp/response.go`
- Create: `glhttp/client.go`
- Create: `glhttp/http_test.go`

**Interfaces:**
- Consumes: `glerror.CodeError` is used by `glhttp.Fail`.
- Produces:
  - `glerror.CodeError`
  - `glerror.New(code int, message string) *CodeError`
  - `glerror.Wrap(code int, message string, err error) *CodeError`
  - `glerror.IsCode(err error, code int) bool`
  - `glerror.From(err error) (*CodeError, bool)`
  - `glhttp.Response[T any]`
  - `glhttp.Success[T any](data T) Response[T]`
  - `glhttp.Fail(code int, message string) Response[any]`
  - `glhttp.Client`
  - `glhttp.NewClient(timeout time.Duration) *Client`
  - `(*Client).GetJSON(ctx context.Context, url string, headers map[string]string, out any) error`
  - `(*Client).PostJSON(ctx context.Context, url string, headers map[string]string, body any, out any) error`
  - `(*Client).PostForm(ctx context.Context, url string, headers map[string]string, form url.Values, out any) error`
  - `glhttp.MultipartFile`
  - `(*Client).PostMultipart(ctx context.Context, url string, headers map[string]string, fields url.Values, files []MultipartFile, out any) error`
  - `(*Client).PostFile(ctx context.Context, url string, headers map[string]string, fields url.Values, fieldName string, filePath string, out any) error`

- [ ] **Step 1: Write failing tests for `glerror`**

```go
base := errors.New("db failed")
err := Wrap(50001, "database error", base)
if err.Code != 50001 {
	t.Fatalf("Code = %d", err.Code)
}
if !errors.Is(err, base) {
	t.Fatal("wrapped error should match base")
}
if !IsCode(err, 50001) {
	t.Fatal("IsCode should be true")
}
```

Run: `go test ./glerror`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 2: Implement `glerror`**

`CodeError` has fields `Code int`, `Message string`, and `Err error`. Implement `Error() string` and `Unwrap() error`. `From` uses `errors.As`.

- [ ] **Step 3: Write failing tests for HTTP responses**

Test success and failure response models:

```go
resp := Success(map[string]string{"id": "1"})
if resp.Code != 0 || resp.Message != "success" {
	t.Fatalf("unexpected response: %+v", resp)
}
```

Run: `go test ./glhttp`

Expected before implementation: FAIL.

- [ ] **Step 4: Write failing tests for HTTP client**

Use `httptest.Server` for GET and POST JSON:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"name":"gltools"}`))
}))
defer server.Close()

var out struct {
	Name string `json:"name"`
}
err := NewClient(time.Second).GetJSON(context.Background(), server.URL, nil, &out)
if err != nil {
	t.Fatal(err)
}
if out.Name != "gltools" {
	t.Fatalf("Name = %q", out.Name)
}
```

Run: `go test ./glhttp`

Expected before implementation: FAIL.

- [ ] **Step 5: Implement `glhttp`**

`Response[T]` uses fields `Code int`, `Message string`, and `Data T`. `Success` uses code `0` and message `"success"`. `Client` wraps `http.Client`, encodes POST body as JSON, applies headers, checks 2xx status, and decodes JSON into `out` when `out != nil`.

- [ ] **Step 6: Run package tests**

Run: `go test ./glerror ./glhttp`

Expected: PASS.

- [ ] **Step 7: Run all current tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add glerror glhttp
git commit -m "feat: add error and http helpers"
```

---

### Task 5: Logging and Configuration Packages

**Files:**
- Create: `gllog/log.go`
- Create: `gllog/log_test.go`
- Create: `glconfig/config.go`
- Create: `glconfig/config_test.go`

**Interfaces:**
- Consumes: dependencies from Task 1.
- Produces:
  - `gllog.Format`
  - `gllog.Config`
  - `gllog.New(cfg Config) (*slog.Logger, error)`
  - `gllog.SetDefault(logger *slog.Logger)`
  - `glconfig.LoadJSON(path string, out any) error`
  - `glconfig.LoadYAML(path string, out any) error`
  - `glconfig.Env(key string, def string) string`
  - `glconfig.EnvInt(key string, def int) int`
  - `glconfig.EnvBool(key string, def bool) bool`

- [ ] **Step 1: Write failing tests for `gllog`**

Test JSON and text logger creation using a `bytes.Buffer`:

```go
buf := new(bytes.Buffer)
logger, err := New(Config{Output: buf, Format: FormatJSON, Level: slog.LevelInfo})
if err != nil {
	t.Fatal(err)
}
logger.Info("hello", slog.String("name", "gltools"))
if !strings.Contains(buf.String(), `"msg":"hello"`) {
	t.Fatalf("log output = %s", buf.String())
}
```

Run: `go test ./gllog`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 2: Implement `gllog`**

Define `type Format string` with constants `FormatJSON` and `FormatText`. `Config` includes `Output io.Writer`, `Format Format`, `Level slog.Level`, and `AddSource bool`. Default output is `os.Stdout`. Return an error for unsupported format.

- [ ] **Step 3: Write failing tests for `glconfig`**

Use `t.TempDir()` to write JSON and YAML files. Use `t.Setenv` for environment tests:

```go
t.Setenv("GLTOOLS_PORT", "8080")
if got := EnvInt("GLTOOLS_PORT", 80); got != 8080 {
	t.Fatalf("EnvInt() = %d", got)
}
```

Run: `go test ./glconfig`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 4: Implement `glconfig`**

Use `github.com/json-iterator/go` with `ConfigCompatibleWithStandardLibrary` for JSON, `gopkg.in/yaml.v3` for YAML, `os.LookupEnv`, `strconv.Atoi`, and `strconv.ParseBool`. `EnvInt` and `EnvBool` return defaults when parsing fails.

- [ ] **Step 5: Run package tests**

Run: `go test ./gllog ./glconfig`

Expected: PASS.

- [ ] **Step 6: Run all current tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add gllog glconfig
git commit -m "feat: add logging and config helpers"
```

---

### Task 6: Retry, Rate Limit, and ID Packages

**Files:**
- Create: `glretry/retry.go`
- Create: `glretry/retry_test.go`
- Create: `gllimit/limit.go`
- Create: `gllimit/limit_test.go`
- Create: `glid/id.go`
- Create: `glid/id_test.go`

**Interfaces:**
- Consumes: dependencies from Task 1.
- Produces:
  - `glretry.Operation`
  - `glretry.Options`
  - `glretry.Do(ctx context.Context, opts Options, operation Operation) error`
  - `glretry.FixedDelay(attempts int, delay time.Duration) Options`
  - `glretry.ExponentialBackoff(attempts int, baseDelay time.Duration, maxDelay time.Duration) Options`
  - `gllimit.Limiter`
  - `gllimit.New(ratePerSecond float64, burst int) *Limiter`
  - `(*Limiter).Allow() bool`
  - `(*Limiter).Wait(ctx context.Context) error`
  - `glid.Generator`
  - `glid.NewGenerator(machineID int64) (*Generator, error)`
  - `(*Generator).Next() (int64, error)`
  - `(*Generator).NextString() (string, error)`

- [ ] **Step 1: Write failing tests for `glretry`**

Test success after retries and context cancellation:

```go
attempts := 0
err := Do(context.Background(), FixedDelay(3, time.Millisecond), func(ctx context.Context) error {
	attempts++
	if attempts < 3 {
		return errors.New("temporary")
	}
	return nil
})
if err != nil {
	t.Fatal(err)
}
if attempts != 3 {
	t.Fatalf("attempts = %d, want 3", attempts)
}
```

Run: `go test ./glretry`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 2: Implement `glretry`**

`Operation` is `func(context.Context) error`. Validate attempts greater than zero. For each failed attempt, wait according to options unless it was the final attempt. Stop immediately when context is canceled.

- [ ] **Step 3: Write failing tests for `gllimit`**

Test burst allowance and context cancellation:

```go
limiter := New(1, 1)
if !limiter.Allow() {
	t.Fatal("first request should be allowed")
}
if limiter.Allow() {
	t.Fatal("second immediate request should be rejected")
}
```

Run: `go test ./gllimit`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 4: Implement `gllimit`**

Wrap `golang.org/x/time/rate.Limiter`. `New` converts `ratePerSecond` to `rate.Limit`. Provide `Allow` and `Wait` methods with Chinese GoDoc.

- [ ] **Step 5: Write failing tests for `glid`**

Test monotonic IDs and invalid machine ID:

```go
gen, err := NewGenerator(1)
if err != nil {
	t.Fatal(err)
}
first, err := gen.Next()
if err != nil {
	t.Fatal(err)
}
second, err := gen.Next()
if err != nil {
	t.Fatal(err)
}
if second <= first {
	t.Fatalf("ids should increase: %d <= %d", second, first)
}
```

Run: `go test ./glid`

Expected before implementation: FAIL. Expected after implementation: PASS.

- [ ] **Step 6: Implement `glid`**

Use a Snowflake-style layout with millisecond timestamp, 10-bit machine ID, and 12-bit sequence. Use a mutex to make `Generator` safe for concurrent use. Return an error when machine ID is outside `0..1023` or the system clock moves backward.

- [ ] **Step 7: Run package tests**

Run: `go test ./glretry ./gllimit ./glid`

Expected: PASS.

- [ ] **Step 8: Run all current tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 9: Commit**

```bash
git add glretry gllimit glid
git commit -m "feat: add retry limit and id helpers"
```

---

### Task 7: Usage Documentation and Final Verification

**Files:**
- Modify: `README.md`
- Create: `docs/usage.md`

**Interfaces:**
- Consumes: all public APIs from Tasks 2-6.
- Produces: Chinese quick-start README and complete Chinese usage guide.

- [ ] **Step 1: Expand README**

README must include:

```markdown
## 包列表

| 包名 | 说明 |
| --- | --- |
| glstrings | 字符串工具 |
| gltime | 时间工具 |
| glconv | 类型转换工具 |
| glslice | 切片工具 |
| glmap | Map 工具 |
| glrand | 安全随机工具 |
| glcrypto | 摘要、签名和 Base64 工具 |
| gljson | JSON 工具 |
| glfile | 文件和路径工具 |
| glhttp | HTTP 响应和客户端工具 |
| gllog | slog 日志工具 |
| glconfig | 配置读取工具 |
| glerror | 错误码工具 |
| glretry | 重试工具 |
| gllimit | 限流工具 |
| glid | 分布式 ID 工具 |
```

- [ ] **Step 2: Create `docs/usage.md`**

The document must contain one section for each package. Each section includes purpose, import path, common functions, and one runnable code example. Use this format for each package:

````markdown
## glstrings

`glstrings` 提供字符串判空、截断和命名转换等能力。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glstrings"
)

func main() {
	fmt.Println(glstrings.IsBlank("   "))
}
```
````

- [ ] **Step 3: Check exported API comments and vet output**

Run: `go vet ./...`

Expected: PASS. Then inspect every exported identifier in the changed packages and confirm it has a Chinese GoDoc comment before continuing.

- [ ] **Step 4: Run full tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 5: Run final vet**

Run: `go vet ./...`

Expected: PASS.

- [ ] **Step 6: Check repository status**

Run: `git status --short`

Expected: only intended documentation changes before commit.

- [ ] **Step 7: Commit**

```bash
git add README.md docs/usage.md
git commit -m "docs: add toolkit usage guide"
```

---

## Self-Review Checklist

- Every package in the design spec has a task.
- Every task has a RED test step before implementation.
- Every public API listed in interfaces has a corresponding package and test responsibility.
- Documentation tasks include README and `docs/usage.md`.
- Final verification includes `go test ./...` and `go vet ./...`.
- No task introduces framework binding, database clients, Redis clients, JWT, object storage, dependency injection, or CLI behavior.
