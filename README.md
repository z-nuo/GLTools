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
| gllog | slog 日志工具 |
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

## 使用说明

完整说明见 `docs/usage.md`。
