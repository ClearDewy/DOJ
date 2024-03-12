/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package config

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type LogrusFormatter struct{}

func (f *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	var levelColor func(...interface{}) string
	switch entry.Level {
	case logrus.ErrorLevel:
		levelColor = color.New(color.FgRed).SprintFunc()
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow).SprintFunc()
	case logrus.InfoLevel:
		levelColor = color.New(color.FgGreen).SprintFunc()
	case logrus.DebugLevel:
		levelColor = color.New(color.FgWhite).SprintFunc()
	default:
		levelColor = color.New(color.FgWhite).SprintFunc()
	}

	// Split the message into lines and add padding to all lines except the first
	lines := strings.Split(message, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = fmt.Sprintf("%-7s\t%s\t%s", "", "", lines[i]) // Adding padding
	}
	// Join the lines back together
	message = strings.Join(lines, "\n")

	return []byte(fmt.Sprintf("%-7s\t%s\t%s\n", levelColor(level), timestamp, message)), nil
}

// InitLogrus 初始化日志格式
func InitLogrus() {
	logrus.SetFormatter(new(LogrusFormatter))
}
