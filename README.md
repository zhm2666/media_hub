# 镜像生成（应用打包）
## 镜像打包
```shell
docker build -t 2410/mediahub:1.0.0 -t 192.168.239.161:5000/2410/mediahub:1.0.0 .
```
## 镜像推送的注册中心
```shell
docker push 192.168.239.161:5000/2410/mediahub:1.0.0
```

# 部署服务
## 创建配置
``` shell
docker config create 2410-mediahub-conf dev.config.yaml
```
## 启动服务
``` shell
docker service create --name 2410-mediahub \
-p 8080:8080 \
--config src=2410-mediahub-conf,target=/app/config.yaml \
--replicas 2 \
--health-cmd "curl -f http://localhost:8080/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
--network 2410-shorturl-net \
192.168.239.161:5000/2410/mediahub:1.0.0
```