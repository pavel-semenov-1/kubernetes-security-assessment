FROM node:20

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . ./
COPY entrypoint.sh /app/entrypoint.sh
RUN npm run build && \
    chmod +x /app/entrypoint.sh

EXPOSE 3000

ENTRYPOINT /app/entrypoint.sh