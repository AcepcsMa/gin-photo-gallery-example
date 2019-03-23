#!/bin/sh

#echo $DB_HOST
#echo $REDIS_HOST

sed -i "s/<SERVER_PORT>/$SERVER_PORT/g" conf/server.conf
sed -i "s/<DB_HOST>/$DB_HOST/g" conf/server.conf
sed -i "s/<REDIS_HOST>/$REDIS_HOST/g" conf/server.conf
sed -i "s/<ES_HOST>/$ES_HOST/g" conf/server.conf

#cat conf/server.conf

exec "$@"
