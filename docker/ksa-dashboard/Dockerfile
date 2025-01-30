FROM node:20 AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install --frozen-lockfile

COPY . .

RUN npm run build

FROM node:20-slim

WORKDIR /app

COPY --from=builder --chown=node:node /app ./

USER node

ENV PORT 3000

EXPOSE $PORT

CMD ["npm", "run", "start"]
