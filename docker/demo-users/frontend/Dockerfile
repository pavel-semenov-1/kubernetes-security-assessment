# use the latest image from Dockerhub
# use image tag instead of hash
FROM node:latest

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . ./
COPY entrypoint.sh /app/entrypoint.sh
RUN npm run build && \
    chmod +x /app/entrypoint.sh

# sensitive information in environment variables
ENV AUTH_TOKEN=QmVhcmVyIHh3SmRhMDl4T0NhNWFQUXpxM2NjeHF3Vwo=

EXPOSE 3000
# expose unnecessary ports
EXPOSE 3000 22

# run as root
USER root

ENTRYPOINT /app/entrypoint.sh