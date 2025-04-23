# 构建阶段
# FROM registry.cn-hangzhou.aliyuncs.com/mszlu/go:1.19 AS build-stage

# WORKDIR /app
# COPY . .

# RUN go env -w GO111MODULE=on
# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN apt-get update && apt-get install -y git
# RUN go mod download

# # 调试：列出 ucenter-api 目录内容
# RUN ls -la /app/ucenter-api

# # 构建二进制文件
# RUN CGO_ENABLED=0 GOOS=linux go build -o /user-api ./ucenter-api

# # 运行阶段
# FROM registry.cn-hangzhou.aliyuncs.com/mszlu-gcrio/distroless_base-debian11 AS build-release-stage

# WORKDIR /
# COPY --from=build-stage /app/user-api /user-api
# COPY --from=build-stage /app/ucenter-api/etc/conf.yaml /etc/conf.yaml

# CMD ["/user-api"]

# 我们需要引入到基础容器
FROM alpine:3.18

# 注意看这里我们的写法, 在sh中 && 可以表示下一条命令连续执行，而 \ 则是命令的分隔符号
# 思考：为什么这么写？
RUN mkdir /go && cd /go \
    && wget --no-check-certificate https://golang.google.cn/dl/go1.21.0.linux-amd64.tar.gz \
    && tar -C /usr/local -zxf go1.21.0.linux-amd64.tar.gz \
    && rm -rf /go/go1.21.0.linux-amd64.tar.gz \
    && mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 
   
# 配置系统环境变量
# 在dockerfile中对容器的系统环境变量配置统一采用ENV这个关键词定义
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .

# RUN go mod download
RUN apk add --no-cache git ca-certificates

# 这个命令可以让我们的docker容器在启动的时候就执行下面的命令
# 把原本在docker run中的命令放到dockerfile中，并示意启动容器的时候执行
# 但是如果在docker run后跟新的命令会代替CMD中的命令
CMD ["ping", "www.baidu.com"]


# FROM alpine:3.18

# # 创建工作目录
# RUN mkdir /go/user-api && cd /go

# # 安装必要的工具
# RUN apk add --no-cache wget ca-certificates git

# # 下载并安装 Golang
# RUN wget --no-check-certificate https://golang.google.cn/dl/go1.21.0.linux-amd64.tar.gz -O /go/go1.21.0.linux-amd64.tar.gz
# RUN tar -C /go -xzf /go/go1.21.0.linux-amd64.tar.gz

# # 设置 Go 环境变量
# ENV GOROOT=/go/go1.21.0
# ENV PATH=$GOROOT/bin:$PATH

# # 设置 Go 的代理
# RUN go env -w GOPROXY=https://goproxy.cn,direct

# # 复制 go.mod 和 go.sum 到容器
# COPY go.mod go.sum /go/user-api/

# # 下载 Go 模块
# RUN go mod download


# CMD ["ping", "www.baidu.com"]

