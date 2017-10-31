FROM alpine:3.6

ADD bin/imagelist_linux_amd64 /bin/imagelist

ENTRYPOINT [ "/bin/imagelist" ]
