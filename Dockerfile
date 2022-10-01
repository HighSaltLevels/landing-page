FROM node:17-alpine3.12

WORKDIR /landing-page

RUN npm install http-server -g && \
    touch /landing-page/health # For a health check endpoint

COPY assets /landing-page/assets
COPY images /landing-page/images
COPY favicon.ico index.html 404.html /landing-page/

CMD npx http-server --cors