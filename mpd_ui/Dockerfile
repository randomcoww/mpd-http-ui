FROM node:10-alpine

WORKDIR /usr/src/app
ADD . /usr/src/app

RUN set -x \
  \
  && npm install -g -s --no-progress yarn \
  && yarn \
  && yarn run build \
  && yarn cache clean

CMD [ "yarn", "start" ]
