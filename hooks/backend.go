package hooks

import "github.com/sirupsen/logrus"

type Backend interface {
	Log(level logrus.Level, message []byte) error
}

/*
When logrus writes message to itself output write, it will call formatter to
format event. And when call hooks, we can not use previous formatted output.
So when we want to ignore itself output, we can set logrus will a NullFormatter,
and set hook with formatter we want, to avoid useless cost.
*/
type backendHook struct {
	backend   Backend
	formatter logrus.Formatter
	levels    []logrus.Level
}

func NewBackendHook(
	backend Backend, formatter logrus.Formatter, levels []logrus.Level,
) *backendHook {
	return &backendHook{
		backend:   backend,
		formatter: formatter,
		levels:    levels,
	}
}

func (hook *backendHook) Fire(entry *logrus.Entry) error {
	message, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	if err = hook.backend.Log(entry.Level, message); err != nil {
		return err
	}
	return nil
}

func (hook *backendHook) Levels() []logrus.Level {
	return hook.levels
}