/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package utils

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

func HandleError(err error, msg string) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			logrus.Errorf("%s.At %s:%d. \n%v", msg, file, line, err)
		} else {
			logrus.Errorf("%s.Can't get the location. \n%v", msg, err)
		}
	}
}
