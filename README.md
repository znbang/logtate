# logtate

logtate is a simple logger with log rotation.

## Example

```go
import "github.com/znbang/logtate"

func main() {
    logger := logtate.New(logtate.Option{
		Path: "hello.log",
		MaxBackup: 5, // 5 backups
		MaxSize: 5, // MB
	})
	
	log.SetOutput(logger)
	log.Println("hello")
}
```