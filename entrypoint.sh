#!/bin/sh
[ "$DB_URL" ] || DB_URL='user:password@tcp(192.168.10.100:3306)/kingfisher'
[ "$LISTEN" ] || LISTEN=0.0.0.0
[ "$PORT" ] || PORT=8080
[ "$RABBITMQ_URL" ] || RABBITMQ_URL='amqp://user:password@king-rabbitmq:5672/'

/usr/local/bin/king-inspect -dbURL=$DB_URL -listen=$LISTEN:$PORT -rabbitMQURL=$RABBITMQ_URL
