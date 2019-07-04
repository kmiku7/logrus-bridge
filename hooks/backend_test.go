package hooks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	tmpFile, err := ioutil.TempFile("", "test_backend")
	if err != nil {
		t.Errorf("Open temporary file failed, err: %v", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	ioBackend := &IOBackend{
		Writer: tmpFile,
	}
	hook := NewBackendHook(
		ioBackend,
		&logrus.TextFormatter{},
		logrus.AllLevels,
	)
	logger := logrus.New()
	logger.AddHook(hook)

	infoMessage := "This is a info message."
	logger.Info(infoMessage)

	warningMessage := "This is a warning message."
	logger.Warn(warningMessage)

	if _, err = tmpFile.Seek(0, 0); err != nil {
		t.Errorf("Seek temporary file failed, err: %v", err)
	}
	contents, err := ioutil.ReadAll(tmpFile)
	if err != nil {
		t.Errorf("Read result contents failed, err: %v", err)
	}
	if !bytes.Contains(contents, []byte("level=info msg=\"This is a info message.\"")) {
		t.Errorf("Info message does not match, output: %s", string(contents))
	}
	if !bytes.Contains(contents, []byte("level=warning msg=\"This is a warning message.\"")) {
		t.Errorf("Warn message does not match, output: %s", string(contents))
	}
}