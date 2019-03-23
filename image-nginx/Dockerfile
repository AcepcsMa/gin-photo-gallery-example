FROM nginx:1.14-alpine

COPY server.crt server.key bootstrap.sh /

CMD ["nginx", "-c", "/etc/nginx/conf.d/my.conf", "-g", "daemon off;"]

ENTRYPOINT ["./bootstrap.sh"]
