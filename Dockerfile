FROM node:17-alpine3.12

WORKDIR /landing-page

RUN npm install http-server -g

COPY assets /landing-page/assets
COPY images /landing-page/images
COPY index.html /landing-page/index.html

CMD npx http-server
