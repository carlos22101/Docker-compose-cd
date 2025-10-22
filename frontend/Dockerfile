FROM node:20-slim

WORKDIR /usr/src/app

COPY package.json ./
RUN npm install --silent

COPY server.js ./
COPY public ./public

EXPOSE 80
CMD ["node", "server.js"]
