# 基础镜像
FROM golang:1.19-alpine3.15
# 设置go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \ 
    GOPROXY="https://goproxy.cn"

# 切换工作目录
WORKDIR "/build"
# 拷贝代码到容器中
COPY . .
# go构建可执行文件
RUN go build -o my_todo ./main.go

# 暴露端口
EXPOSE 8888
# 启动容器时候运行的命令
ENTRYPOINT ["./my_todo"]