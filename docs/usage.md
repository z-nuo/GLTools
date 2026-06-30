# GLTools 使用指南

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

`glhttp` 提供统一 JSON 响应结构和简单 JSON HTTP 客户端。

导入路径：`github.com/z-nuo/GLTools/glhttp`

常用类型和函数：`Response`、`Success`、`Fail`、`Client`、`NewClient`、`GetJSON`、`PostJSON`。

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/z-nuo/GLTools/glhttp"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "ok"})
	}))
	defer server.Close()

	var out map[string]string
	client := glhttp.NewClient(time.Second)
	if err := client.GetJSON(context.Background(), server.URL, nil, &out); err != nil {
		panic(err)
	}
	fmt.Println(glhttp.Success(out).Data["message"])
}
```

## gllog

`gllog` 提供基于 `log/slog` 的日志器创建和默认日志器设置能力。

导入路径：`github.com/z-nuo/GLTools/gllog`

常用类型和函数：`Format`、`FormatJSON`、`FormatText`、`Config`、`New`、`SetDefault`。

```go
package main

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/z-nuo/GLTools/gllog"
)

func main() {
	var buf bytes.Buffer
	logger, err := gllog.New(gllog.Config{
		Output: &buf,
		Format: gllog.FormatText,
		Level:  slog.LevelInfo,
	})
	if err != nil {
		panic(err)
	}
	logger.Info("hello")
	fmt.Println(buf.Len() > 0)
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
