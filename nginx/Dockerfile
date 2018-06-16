FROM nginx:1.14

LABEL maintainer="Ivan de la Beldad Fernandez <ivandelabeldad@gmail.com>"

RUN rm -f /etc/nginx/conf.d/default.conf

COPY docker /docker

CMD [ "sh", "/docker/entrypoint.sh" ]
