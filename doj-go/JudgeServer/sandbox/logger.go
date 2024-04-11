/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sandbox

import (
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/env"
)

type Logger struct {
	env.Logger
}

func (*Logger) Debug(args ...any) {
	logrus.Debug(args)
}
func (*Logger) Info(args ...any) {
	logrus.Info(args)
}
func (*Logger) Warn(args ...any) {
	logrus.Warn(args)
}
func (*Logger) Error(args ...any) {
	logrus.Error(args)
}
