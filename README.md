# GLTools

GLTools 是一个面向 Go 后端服务的通用工具库，封装字符串、时间、类型转换、切片、Map、随机数、加密摘要、JSON、文件、HTTP、日志、配置、错误码、重试、限流和 ID 生成等常用能力。

## 安装

```bash
go get github.com/z-nuo/GLTools
```

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
| gllog | zap 日志工具 |
| glconfig | 配置读取工具 |
| glerror | 错误码工具 |
| glretry | 重试工具 |
| gllimit | 限流工具 |
| glid | 分布式 ID 工具 |

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

## 使用指南

本文档按包介绍 GLTools 的常用能力。示例均为可直接运行的最小程序。

## glstrings

`glstrings` 提供字符串判空、截断和命名转换等能力。

导入路径：`github.com/z-nuo/GLTools/glstrings`

常用函数：`IsBlank`、`Trim`、`Truncate`、`SnakeToCamel`、`CamelToSnake`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glstrings"
)

func main() {
	fmt.Println(glstrings.IsBlank("   "))
	fmt.Println(glstrings.Truncate("你好世界", 2))
	fmt.Println(glstrings.SnakeToCamel("user_name"))
}
```

## gltime

`gltime` 提供时间格式化、解析、毫秒时间戳和日期边界计算能力。

导入路径：`github.com/z-nuo/GLTools/gltime`

常用函数：`Format`、`Parse`、`UnixMilli`、`FromUnixMilli`、`DayStart`、`DayEnd`。

```go
package main

import (
	"fmt"
	"time"

	"github.com/z-nuo/GLTools/gltime"
)

func main() {
	now := time.Date(2026, 6, 30, 15, 4, 5, 0, time.Local)
	fmt.Println(gltime.Format(now, "2006-01-02"))
	fmt.Println(gltime.DayStart(now).Format(time.DateTime))
}
```

## glconv

`glconv` 提供字符串到基础类型的转换和任意值字符串化能力。

导入路径：`github.com/z-nuo/GLTools/glconv`

常用函数：`ToInt`、`ToIntDefault`、`ToInt64`、`ToFloat64`、`ToBool`、`String`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glconv"
)

func main() {
	port := glconv.ToIntDefault("8080", 80)
	enabled, _ := glconv.ToBool("true")
	fmt.Println(port, enabled, glconv.String(123))
}
```

## glslice

`glslice` 提供切片包含判断、去重、过滤和字符串清理能力。

导入路径：`github.com/z-nuo/GLTools/glslice`

常用函数：`Contains`、`Unique`、`Filter`、`CompactStrings`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glslice"
)

func main() {
	items := []int{1, 2, 2, 3}
	fmt.Println(glslice.Contains(items, 2))
	fmt.Println(glslice.Unique(items))
	fmt.Println(glslice.Filter(items, func(v int) bool { return v > 1 }))
}
```

## glmap

`glmap` 提供 Map 键值提取、键存在判断和合并能力。

导入路径：`github.com/z-nuo/GLTools/glmap`

常用函数：`Keys`、`Values`、`HasKey`、`Merge`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glmap"
)

func main() {
	left := map[string]int{"a": 1}
	right := map[string]int{"b": 2, "a": 3}
	merged := glmap.Merge(left, right)
	fmt.Println(glmap.HasKey(merged, "a"))
	fmt.Println(merged["a"])
}
```

## glrand

`glrand` 提供基于 `crypto/rand` 的安全随机数字码和字符串生成能力。

导入路径：`github.com/z-nuo/GLTools/glrand`

常用函数：`NumericCode`、`String`、`StringFromCharset`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glrand"
)

func main() {
	code, err := glrand.NumericCode(6)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(code))
}
```

## glcrypto

`glcrypto` 提供 MD5、SHA256、HMAC-SHA256 摘要和 Base64 编解码能力。

导入路径：`github.com/z-nuo/GLTools/glcrypto`

常用函数：`MD5Hex`、`SHA256Hex`、`HMACSHA256Hex`、`Base64Encode`、`Base64Decode`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glcrypto"
)

func main() {
	encoded := glcrypto.Base64Encode("hello")
	decoded, _ := glcrypto.Base64Decode(encoded)
	fmt.Println(glcrypto.SHA256Hex("hello"))
	fmt.Println(decoded)
}
```

## gljson

`gljson` 提供 JSON 编码、解码、校验和格式化能力。

导入路径：`github.com/z-nuo/GLTools/gljson`

常用函数：`Marshal`、`Unmarshal`、`Valid`、`Pretty`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/gljson"
)

func main() {
	body, _ := gljson.Pretty(map[string]any{"name": "GLTools"})
	fmt.Println(gljson.Valid(body))
	fmt.Println(body)
}
```

## glfile

`glfile` 提供文件存在性判断、目录创建、文本读写和路径处理能力。

导入路径：`github.com/z-nuo/GLTools/glfile`

常用函数：`Exists`、`IsFile`、`IsDir`、`EnsureDir`、`ReadText`、`WriteText`、`Ext`、`Join`。

```go
package main

import (
	"fmt"
	"os"

	"github.com/z-nuo/GLTools/glfile"
)

func main() {
	dir, err := os.MkdirTemp("", "gltools-example-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := glfile.Join(dir, "hello.txt")
	if err := glfile.WriteText(path, "hello"); err != nil {
		panic(err)
	}
	text, _ := glfile.ReadText(path)
	fmt.Println(glfile.Exists(path), glfile.Ext(path), text)
}
```

## glhttp

`glhttp` 提供统一 JSON 响应结构、JSON HTTP 客户端、表单 POST 请求和 multipart 文件上传能力。

导入路径：`github.com/z-nuo/GLTools/glhttp`

常用类型和函数：`Response`、`Success`、`Fail`、`Client`、`NewClient`、`NewClientWithHTTPClient`、`GetJSON`、`PostJSON`、`PostForm`、`MultipartFile`、`PostMultipart`、`PostFile`。

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/z-nuo/GLTools/glhttp"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		if r.Form.Get("name") != "GLTools" {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"message":"ok"}`))
	}))
	defer server.Close()

	var out map[string]string
	client := glhttp.NewClient(time.Second)
	form := url.Values{"name": []string{"GLTools"}}
	if err := client.PostForm(context.Background(), server.URL, nil, form, &out); err != nil {
		panic(err)
	}
	fmt.Println(glhttp.Success(out).Data["message"])
}
```

本地文件上传可以使用 `PostFile`；需要从内存、网络流或多个文件上传时，使用 `MultipartFile` 和 `PostMultipart`。

```text
fields := url.Values{"name": []string{"GLTools"}}
err := client.PostFile(context.Background(), uploadURL, nil, fields, "file", "./demo.txt", &out)
```

## gllog

`gllog` 提供基于 `go.uber.org/zap` 的高性能结构化日志能力，支持控制台/JSON 输出、默认日志器、调用方信息、错误堆栈、链路字段，以及按小时或按天切分日志文件。

导入路径：`github.com/z-nuo/GLTools/gllog`

常用类型和函数：`Format`、`FormatJSON`、`FormatConsole`、`Level`、`LevelInfo`、`Rotate`、`RotateHourly`、`RotateDaily`、`Config`、`New`、`SetDefault`、`L`、`S`、`Sync`、`WithTrace`、`TraceID`、`SpanID`、`TraceFields`、`WithContext`、`DebugContext`、`InfoContext`、`WarnContext`、`ErrorContext`。

企业级日志库通常需要满足这些要求：

- 支持结构化字段，便于日志采集、检索和告警。
- 支持 debug、info、warn、error 等日志级别。
- 支持 JSON 格式，便于接入 ELK、Loki、Datadog 等平台。
- 支持控制台格式，便于本地开发排查。
- 支持调用方文件和行号，便于定位问题。
- 支持 error 级别堆栈，便于线上故障分析。
- 支持按小时或按天切分日志文件，便于归档和清理。
- 支持 stdout、文件、自定义 writer 和同时输出。
- 支持全局默认 logger，降低业务代码接入成本。
- 支持从 `context.Context` 自动写入 `trace_id` 和 `span_id`，便于按请求链路检索日志。
- 提供 Sync 能力，进程退出前可主动刷新日志缓冲。

说明：`gllog` 的链路追踪能力是日志侧的链路字段透传和输出，适合把一次请求、一次 RPC 或一个任务的日志串起来检索；它不替代 OpenTelemetry、Jaeger、SkyWalking 等完整分布式追踪系统。

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"github.com/z-nuo/GLTools/gllog"
)

func main() {
	var buf bytes.Buffer
	logger, err := gllog.New(gllog.Config{
		Output: &buf,
		Format: gllog.FormatJSON,
		Level:  gllog.LevelInfo,
	})
	if err != nil {
		panic(err)
	}
	gllog.SetDefault(logger)

	ctx := gllog.WithTrace(context.Background(), "trace-001", "span-001")
	gllog.InfoContext(ctx, "hello", zap.String("name", "GLTools"))

	fmt.Println(strings.Contains(buf.String(), `"trace_id":"trace-001"`))
}
```

## glconfig

`glconfig` 提供 JSON、YAML 配置文件加载和环境变量读取能力。

导入路径：`github.com/z-nuo/GLTools/glconfig`

常用函数：`LoadJSON`、`LoadYAML`、`Env`、`EnvInt`、`EnvBool`。

```go
package main

import (
	"fmt"
	"os"

	"github.com/z-nuo/GLTools/glconfig"
)

func main() {
	_ = os.Setenv("GLTOOLS_PORT", "8080")
	fmt.Println(glconfig.EnvInt("GLTOOLS_PORT", 80))
	fmt.Println(glconfig.Env("GLTOOLS_ENV", "dev"))
}
```

## glerror

`glerror` 提供带业务错误码的错误创建、包装和识别能力。

导入路径：`github.com/z-nuo/GLTools/glerror`

常用类型和函数：`CodeError`、`New`、`Wrap`、`From`、`IsCode`。

```go
package main

import (
	"errors"
	"fmt"

	"github.com/z-nuo/GLTools/glerror"
)

func main() {
	err := glerror.Wrap(1001, "保存失败", errors.New("database timeout"))
	codeErr, ok := glerror.From(err)
	fmt.Println(ok, codeErr.Code, glerror.IsCode(err, 1001))
}
```

## glretry

`glretry` 提供固定间隔和指数退避的重试执行能力。

导入路径：`github.com/z-nuo/GLTools/glretry`

常用类型和函数：`Operation`、`Options`、`FixedDelay`、`ExponentialBackoff`、`Do`。

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/z-nuo/GLTools/glretry"
)

func main() {
	attempts := 0
	err := glretry.Do(context.Background(), glretry.FixedDelay(3, time.Millisecond), func(context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary")
		}
		return nil
	})
	fmt.Println(err == nil, attempts)
}
```

## gllimit

`gllimit` 提供基于令牌桶的限流能力。

导入路径：`github.com/z-nuo/GLTools/gllimit`

常用类型和函数：`Limiter`、`New`、`Allow`、`Wait`。

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/z-nuo/GLTools/gllimit"
)

func main() {
	limiter := gllimit.New(1, 1)
	fmt.Println(limiter.Allow())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	fmt.Println(limiter.Wait(ctx) != nil)
}
```

## glid

`glid` 提供 Snowflake 风格的分布式递增 ID 生成能力。

导入路径：`github.com/z-nuo/GLTools/glid`

常用类型和函数：`Generator`、`NewGenerator`、`Next`、`NextString`。

```go
package main

import (
	"fmt"

	"github.com/z-nuo/GLTools/glid"
)

func main() {
	generator, err := glid.NewGenerator(1)
	if err != nil {
		panic(err)
	}
	id, _ := generator.NextString()
	fmt.Println(len(id) > 0)
}
```
