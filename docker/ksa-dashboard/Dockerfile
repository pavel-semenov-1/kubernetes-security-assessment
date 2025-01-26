FROM node:20 AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install --frozen-lockfile

COPY . .

RUN npm run build

FROM node:20-slim

WORKDIR /app

COPY --from=builder --chown=node:node /app ./

RUN npm install -g serve

USER node

ENV PORT 3000

EXPOSE $PORT

CMD ["sh", "-c", "serve -s out -l ${PORT}"]
