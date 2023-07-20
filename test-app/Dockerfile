FROM node:alpine
USER 1000:1000
EXPOSE 8080
WORKDIR /app
COPY package.json .
RUN npm install
COPY server.js .
CMD ["node", "server.js"]
