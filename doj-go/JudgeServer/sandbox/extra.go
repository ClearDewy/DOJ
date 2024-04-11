/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sandbox

import (
	"context"
	"github.com/ClearDewy/go-pkg/logrus"
	"github.com/criyle/go-judge/pb"
)

func DeleteFileByFileID(fileIDs []string) {
	for _, value := range fileIDs {
		_, err := Client.FileDelete(context.Background(), &pb.FileID{
			FileID: value,
		})
		if err != nil {
			logrus.ErrorM(err, "删除缓存文件异常")
		}
	}
}
func DeleteFileByResult(result *pb.Response_Result) {
	if result == nil {
		return
	}
	for _, value := range result.FileIDs {
		_, err := Client.FileDelete(context.Background(), &pb.FileID{
			FileID: value,
		})
		if err != nil {
			logrus.Error()
			logrus.ErrorM(err, "删除缓存文件异常")
		}
	}
}
