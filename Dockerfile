FROM alpine:3.10

ENV TIME_ZONE Asia/Shanghai
RUN set -xe \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --no-cache add tzdata \
    && echo "${TIME_ZONE}" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime \
    && mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ADD bin/king-inspect /usr/local/bin

CMD /usr/local/bin/king-inspect -dbURL='user:password@tcp(192.168.10.100:3306)/kingfisher' -listen=0.0.0.0:8080 -rabbitMQURL='amqp://user:password@king-rabbitmq:5672/'

EXPOSE 8080