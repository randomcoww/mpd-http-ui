FROM node:10-alpine as BUILD

WORKDIR /build
COPY . .

RUN set -x \
  \
  && yarn install \
  && NODE_ENV=production yarn build

FROM alpine:edge

RUN set -x \
  \
  && apk add --no-cache \
    nginx \
  \
  && mkdir -p /run/nginx \
  && ln -sf /dev/stdout /var/log/nginx/access.log \
  && ln -sf /dev/stderr /var/log/nginx/error.log

COPY --from=BUILD /build/dist /var/www/dist
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

CMD [ "nginx", "-g", "daemon off;" ]
