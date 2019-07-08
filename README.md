# logrus-bridge
Use logrus and log frontend, and bridge to various storage backend.


**Code Example:**
```golang
import (
  "io/ioutil"

  "github.com/kmiku7/logrus-bridge/formatter"
  "github.com/kmiku7/logrus-bridge/hooks"
  "github.com/kmiku7/logrus-bridge/logger"
  "github.com/sirupsen/logrus"
)

var GlobalLoggerClient logger.Logger


func InitGlobalLogger() {
  var backend hooks.Backend // = SomePackage.NewBackend(...)

  logClient := logrus.New()
  logClient.Out = ioutil.Discard
  logClient.Formatter = formatter.EmptyFormatter(0)
  logClient.SetLevel(logrus.DebugLevel)
  
  hook := hooks.NewBackendHook(
    backend,
    &logrus.TextFormatter{},
    logrus.AllLevels)
  logClient.AddHook(hook)
  
  GlobalLoggerClient = logClient
}

func main() {
  InitGlobalLogger()
  
  GlobalLoggerClient.Info("Hello World!")
}
```
