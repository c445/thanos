FROM alpine:latest
LABEL maintainer="The Thanos Authors"

# add capabilities to allow running ulimit -l inside the container
RUN apk update && apk add libcap && setcap cap_sys_resource,cap_ipc_lock=+ep /bin/busybox

COPY /thanos_tmp_for_docker /bin/thanos

ENTRYPOINT [ "/bin/thanos" ]
