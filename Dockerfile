FROM 5eplay-registry.cn-hangzhou.cr.aliyuncs.com/yunao/go:1.18.1-alpine-git-grcp-health AS builder

ARG BUILD_DIR
# 设置当前工作区
WORKDIR /build

#创建app-runner用户, -D表示无密码
RUN adduser -u 10001 -D app-runner

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

COPY go.mod .
COPY go.sum .
RUN go mod download

# 把全部文件添加到/build目录
COPY . .

# 编译: 把main.go编译为可执行的二进制文件, 并命名为app
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w"  -a -o app ${BUILD_DIR}


FROM 5eplay-registry.cn-hangzhou.cr.aliyuncs.com/yunao/alpine:3.10 AS final

WORKDIR /app

COPY --from=builder /build/app /app/
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN mkdir -p /app/logs && chown app-runner /app/logs && chmod 777 /app/logs

USER app-runner

EXPOSE 31301
ENTRYPOINT ["/app/app"]