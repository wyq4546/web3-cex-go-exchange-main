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

# # 我们需要引入到基础容器
# FROM golang:1.20-alpine

# # 安装必要的工具
# RUN apk add --no-cache curl wget vim busybox-extras

# WORKDIR /app

# COPY . .

# # 设置时区
# RUN apk add --no-cache tzdata && \
#     cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
#     echo "Asia/Shanghai" > /etc/timezone

# # 设置环境变量
# ENV GO111MODULE=on \
#     GOPROXY=https://goproxy.cn,direct

# # 暴露端口
# EXPOSE 8081 8888

# CMD ["sh"]


FROM alpine:3.18

RUN apk add --no-cache git ca-certificates

# 注意看这里我们的写法, 在sh中 && 可以表示下一条命令连续执行，而 \ 则是命令的分隔符号
# 思考：为什么这么写？
COPY go1.21.0.linux-amd64.tar.gz /go/
RUN tar -C /usr/local -zxf /go/go1.21.0.linux-amd64.tar.gz \
    && rm -rf /go/go1.21.0.linux-amd64.tar.gz \
    && mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 
   
# 配置系统环境变量
# 在dockerfile中对容器的系统环境变量配置统一采用ENV这个关键词定义
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH

# 设置 Go 的代理
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct


# 安装必要的工具
RUN apk add --no-cache wget ca-certificates git curl vim busybox-extras


# 暴露端口
EXPOSE 8081 8888




