//go:build !linux && !darwin

package grpcexecutor

import (
	"doj-go/JudgeServer/pb"
	"os"
)

func setWinsize(f *os.File, i *pb.StreamRequest_ExecResize) error {
	return nil
}
