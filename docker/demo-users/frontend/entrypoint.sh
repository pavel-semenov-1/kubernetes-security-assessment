#!/bin/sh

sed -i "s|API_BASE_URL=.*$|API_BASE_URL=${API_BASE_URL}|g" .env

HOSTNAME="0.0.0.0" PORT=3000 npm run start