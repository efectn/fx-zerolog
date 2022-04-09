# Zerolog Adapter for Fx

[![Go Reference](https://pkg.go.dev/badge/github.com/efectn/fx-zerolog.svg)](https://pkg.go.dev/github.com/efectn/fx-zerolog)

Zerolog adapter for uber-go/fx/fxevent. 

### Supported Go Versions
- 1.16
- 1.17
- 1.18

### Install

```shell
go get -u github.com/efectn/fx-zerolog@latest
```

### Example
```go
import (
    "github.com/rs/zerolog"
    "go.uber.org/fx"
    "github.com/efectn/fx-zerolog"
)

// ...

func main() {
    logger := zerolog.New(...)

    fx.New(
    	fx.Provide(
    		NewLogger,
    		NewConfig,
    		NewRouter,
    	),
    	fx.Invoke(Listen),

    	WithLogger(
           fxzerolog.Init(logger),
       ),
    )
}
```

### License

fx-zerolog is licensed under the terms of the **GPL-3 License** (see [LICENSE](LICENSE)).
