# Go Logger
Logger for GoLang applications with support of log levels and file output. It can be extremely useful if you don't need very complex libs like `go-logger`, `google/logger`, etc. It has a pretty straight syntax, as well as capabilities.

## Instanciate a Logger

To import it into the project just include it into import section:

```go
import "github.com/dveoz/go/logger/v3"
```

and initiate it in your main func:

```go
// set log level for all other packages
logger.SetLogger("DEBUG")
```

By initiating you will set up logging subsystem ready to send messages to STDOUT only. If that is not enough
you may want to connect file handler by using another function:
```go
// in adddition attach the file handler
logger.SetFileHandler("logs/application.log")
```
Please note that path, sending as a parameter, is **relative** to your working directory.

That's all after that you can start using it. Here are some examples for all levels:

```go
logger.Info("HTTPS Server started on port - %v", serverPort)
```

will print out into both - file and `stdout` the following message:

```
2019/12/24 14:53:34.346930 Server.go:37: [INFO   ] HTTPS Server started on port - 9000
```
