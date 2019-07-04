package hooks

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

type IOBackend struct {
	Writer io.Writer
}

func (b *IOBackend) Log(level logrus.Level, message []byte) error {
	levelName, err := level.MarshalText()
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(b.Writer, "%s %s", levelName, message); err != nil {
		return err
	}
	return nil
}

func TestBackend(t *testing.T) {
	log := logrus.New()

	tmpFile, err := ioutil.TempFile("", "test_backend")
	if err != nil {
		t.Errorf("Open temporary file failed, err: %v", err)
	}
	defer func() {
		tmpFile.Close()
	}()

	ioBackend := IOBackend{
		Writer: tmpFile,
	}
	hook := NewBackendHook(
		ioBackend,
		logrus.TextFormatter(),
		logrus.AllLevels,
	)
	logger := logrus.New()
	logger.AddHook(hook)

	infoMessage := "This is a info message."
	logger.Info(infoMessage)

	warningMessage := "This is a warning message."
	logger.Warn(warningMessage)

	contents, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Errorf("Read result contents failed, err: %v", err)
	}
	fmt.Println(contents)
}
