FROM quay.io/0voice/golang:1.20 as stage0
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD mediahub /src/mediahub
WORKDIR /src/mediahub
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mediahub .

FROM quay.io/0voice/node:18.16.0 as stage1
RUN npm config set registry https://mirrors.huaweicloud.com/repository/npm/
ADD mediahub-web /src/mediahub-web
WORKDIR /src/mediahub-web
RUN npm install
RUN npm run build

FROM quay.io/0voice/alpine:3.18 as stage2
ADD curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl
MAINTAINER nick
WORKDIR /app
ADD ./mediahub/dev.config.yaml /app/config.yaml
COPY --from=stage0 /src/mediahub/mediahub /app
COPY --from=stage1 /src/mediahub-web/dist /app/www
ENTRYPOINT ["./mediahub"]
CMD ["--config=config.yaml"]