package formatter

import "github.com/sirupsen/logrus"

type EmptyFormatter int

func (_ EmptyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return nil, nil
}