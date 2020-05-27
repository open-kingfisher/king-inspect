#!/bin/sh
export TIME_ZONE="Asia/Shanghai"

sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories 
apk --no-cache add tzdata 
echo "${TIME_ZONE}" > /etc/timezone 
ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime 
mkdir /lib64
ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


[ "$DB_URL" ] || DB_URL='user:password@tcp(192.168.10.100:3306)/kingfisher'
[ "$LISTEN" ] || LISTEN=0.0.0.0
[ "$PORT" ] || PORT=8080
[ "$RABBITMQ_URL" ] || RABBITMQ_URL='amqp://user:password@king-rabbitmq:5672/'

/usr/local/bin/king-inspect -dbURL=$DB_URL -listen=$LISTEN:$PORT -rabbitMQURL=$RABBITMQ_URL
