FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/robot-gitee-keeper-approve
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o robot-gitee-keeper-approve .

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 1000 robot-gitee-keeper-approve && \
    useradd -u 1000 -g robot-gitee-keeper-approve -s /bin/bash -m robot-gitee-keeper-approve

USER robot-gitee-keeper-approve

COPY  --chown=robot-gitee-keeper-approve --from=BUILDER /go/src/github.com/opensourceways/robot-gitee-keeper-approve/robot-gitee-keeper-approve /opt/app/robot-gitee-keeper-approve

ENTRYPOINT ["/opt/app/robot-gitee-keeper-approve"]
