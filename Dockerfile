FROM node:17-alpine3.12

WORKDIR /landing-page

RUN npm install http-server@13.0.2 -g

COPY assets /landing-page/assets
COPY images /landing-page/images
COPY health favicon.ico index.html 404.html /landing-page/

CMD npx http-server --cors --ssl --cert /opt/tls/tls.crt --key /opt/tls/tls.key
