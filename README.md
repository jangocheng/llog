# llog

[![GoDoc](https://godoc.org/github.com/syyongx/llog?status.svg)](https://godoc.org/github.com/syyongx/llog)
[![Go Report Card](https://goreportcard.com/badge/github.com/syyongx/llog)](https://goreportcard.com/report/github.com/syyongx/llog)

llog - Monolog implementation in Go.

## GoDoc

- [godoc for github](https://godoc.org/github.com/syyongx/llog)

## Download & Installs

```
go get github.com/syyongx/llog
```

## Usage

```go
package main

import (
	"github.com/syyongx/llog"
	"github.com/syyongx/llog/handler"
	"github.com/syyongx/llog/formatter"
	"github.com/syyongx/llog/types"
)

func main() {
	logger := NewLogger("my-log")

	file := handler.NewFile("/tmp/llog/go.log", 0664, types.WARNING, true)
	buf := handler.NewBuffer(file, 1, types.WARNING, true)
	f := formatter.NewLine("%Datetime% [%LevelName%] [%Channel%] %Message%\n", time.RFC3339)
	file.SetFormatter(f)

	// push handler
	logger.PushHandler(buf)

	// add log
	logger.Warning("xxx")

	// close and write
	buf.Close()
}
```

## LICENSE

llog source code is licensed under the [MIT](LICENSE) Licence.