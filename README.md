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
