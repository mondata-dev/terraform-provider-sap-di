FROM nginx:1.15.3-alpine


WORKDIR /app

COPY .htpasswd .htpasswd
COPY static static
COPY ./nginx.conf /etc/nginx/nginx.conf

