FROM golang:1.22 AS build

# 创建并设置工作目录
WORKDIR /go/judge

# 复制源代码和配置文件到 /go/judge
COPY jspb $WORKDIR
COPY go.mod $WORKDIR
COPY JudgeServer $WORKDIR

## 下载依赖
#RUN go mod download -x
#
## 编译
#RUN go build -v -o server ./JudgeServer
#
### 移动编译后的文件到 /go/judge
##RUN mkdir -p /go/judge && mv judge-server /go/judge/
##
### 清理
##RUN go clean -modcache
##RUN rm -rf /go/compile
#
## 最终镜像
#FROM golang:1.22
#WORKDIR /go/judge
#COPY --from=build /go/judge/server $WORKDIR
## /go/judge/mount.yaml
#EXPOSE 8088
#
## 设置工作目录和入口点
#ENTRYPOINT ["./server"]
