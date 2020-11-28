ARG SHA="d7342993700f8cd7aba8496c2d0e57be0666e80b4c441925fc6f9361fa81d10e"
FROM alpine@sha256:${SHA}
LABEL maintainer="The Thanos Authors"

# add capabilities to allow running ulimit -l inside the container
RUN apk update && apk add libcap && setcap cap_sys_resource,cap_ipc_lock=+ep /bin/busybox

COPY /thanos_tmp_for_docker /bin/thanos

ENTRYPOINT [ "/bin/thanos" ]
