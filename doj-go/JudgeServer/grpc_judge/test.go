/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package grpc_judge

import (
	"context"
	"doj-go/JudgeServer/sandbox"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/pb"
)

func test() {
	cmd := &pb.Request_CmdType{
		Args: []string{"/usr/bin/g++", "a.cc", "-o", "a"},
		Env:  []string{"PATH=/usr/bin:/bin"},
		Files: []*pb.Request_File{
			{
				File: &pb.Request_File_Local{
					Local: &pb.Request_LocalFile{
						Src: "/doj/judge/problem/1000/0/1.in",
					},
				},
			},
			{
				File: &pb.Request_File_Pipe{
					&pb.Request_PipeCollector{
						Name: "stdout",
						Max:  10240,
					},
				},
			},
			{
				File: &pb.Request_File_Pipe{
					&pb.Request_PipeCollector{
						Name: "stderr",
						Max:  10240,
					},
				},
			},
		},
		CpuTimeLimit: 10000000000,
		MemoryLimit:  104857600,
		ProcLimit:    128,
		CopyIn: map[string]*pb.Request_File{
			"a.cc": {
				File: &pb.Request_File_Memory{
					Memory: &pb.Request_MemoryFile{
						Content: []byte("#include <iostream>\nusing namespace std;\nint main() {\nint a, b;\ncin >> a >> b;\ncout << a + b << endl;\n}"),
					},
				},
			},
		},
		CopyOut: []*pb.Request_CmdCopyOutFile{
			{
				Name: "stdout",
			},
			{
				Name: "stderr",
			},
		},
		CopyOutCached: []*pb.Request_CmdCopyOutFile{
			{
				Name:     "a",
				Optional: true,
			},
		},
	}
	result, err := sandbox.Client.Exec(context.Background(), &pb.Request{
		Cmd: []*pb.Request_CmdType{cmd},
	})
	if err != nil {
		logrus.ErrorM(err, "")
		return
	}
	defer sandbox.DeleteFileByResult(result.Results[0])
	logrus.Info(result)
}
