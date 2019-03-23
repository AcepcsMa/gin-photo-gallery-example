#!/bin/sh

cat > /etc/nginx/conf.d/my.conf << EOF
events {
     worker_connections  1024;
}

http {
server {
    listen       443 ssl;
    server_name  my;

    ssl_certificate      /server.crt;
    ssl_certificate_key  /server.key;

    ssl_session_cache    shared:SSL:1m;
    ssl_session_timeout  5m;

	ssl_ciphers  HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers  on;

    location / {
        #root   html;
        #index  index.html index.htm;
		proxy_pass https://${APP_HOST}:${APP_PORT};
	}
}
}
EOF

exec "$@"
