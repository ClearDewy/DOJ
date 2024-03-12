FROM golang:1.20 AS build

# 创建并设置工作目录
WORKDIR /go/judge

# 复制源代码和配置文件到 /go/compile
COPY . $WORKDIR
COPY ./JudgeServer/sandbox/containerPasswd.txt $WORKDIR
COPY ./JudgeServer/sandbox/mount.yaml $WORKDIR

# 下载依赖
RUN go mod download -x

# 编译
RUN go build -v -o sandbox ./JudgeServer/sandbox/.
# 编译
RUN go build -v -o server ./JudgeServer/server/.

## 移动编译后的文件到 /go/judge
#RUN mkdir -p /go/judge && mv judge-server /go/judge/
#
## 清理
#RUN go clean -modcache
#RUN rm -rf /go/compile

# 最终镜像
FROM golang:1.20
WORKDIR /go/judge
COPY --from=build /go/judge/sandbox /go/judge/server $WORKDIR
# /go/judge/mount.yaml
EXPOSE 5051 8088

# 设置工作目录和入口点
ENTRYPOINT ["./server"]
