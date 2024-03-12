/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_sandbox

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func StartSandbox() chan error {
	errChan := make(chan error)
	// 启动sandbox的命令
	cmd := exec.Command("./sandbox")

	// 创建一个管道，用于读取程序的标准输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		errChan <- err
	}
	// 使用 bufio.Scanner 逐行读取输出
	scanner := bufio.NewScanner(stdoutPipe)
	go func() {
		for scanner.Scan() {
			logrus.Info("Sandbox: ", scanner.Text())
		}
	}()

	go func() {
		errChan <- cmd.Wait()
	}()
	return errChan
}
