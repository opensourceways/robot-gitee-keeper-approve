FROM openeuler/go:1.23.4-oe2403lts as BUILDER
RUN dnf -y install git gcc

ARG USER
ARG PASS
RUN echo "machine github.com login $USER password $PASS" > ~/.netrc

# build binary
WORKDIR /opt/source
COPY . .
RUN go env -w GO111MODULE=on && \
    go env -w CGO_ENABLED=1 && \
    go build -a -o robot-gitee-keeper-approve -buildmode=pie -ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'" .

# copy binary config and utils
FROM openeuler/openeuler:24.03-lts
RUN dnf -y upgrade && \
    dnf in -y shadow && \
    groupadd -g 1000 robot && \
    useradd -u 1000 -g robot -s /bin/bash -m robot

USER robot

COPY --chown=robot --from=BUILDER /opt/source/robot-gitee-keeper-approve /opt/app/robot-gitee-keeper-approve

ENTRYPOINT ["/opt/app/robot-gitee-keeper-approve"]