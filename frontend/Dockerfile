# syntax=docker/dockerfile:1

FROM node:23.3.0-alpine AS dependencies
WORKDIR /app
COPY package.json package-lock.json apply-patches.sh ./
RUN npm ci
RUN apk add --no-cache jq patch && chmod +x apply-patches.sh && ./apply-patches.sh
ARG VERSION=v0.0.0-notset
RUN jq ".version=\"${VERSION}\"" package.json > /tmp/package.json && mv /tmp/package.json package.json

FROM node:23.3.0-alpine AS development
WORKDIR /app
ARG VITE_GEOAPIFY_API_KEY
ENV NODE_ENV=development
COPY --from=dependencies /app ./
COPY eslint.config.js ./
COPY index.html ./
COPY src ./src/
COPY tsconfig.app.json ./
COPY tsconfig.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
ENTRYPOINT ["npm","run","dev"]

FROM development AS builder
ARG VITE_GEOAPIFY_API_KEY
ENV NODE_ENV=production
RUN npm run build

FROM nginx:1.27 AS production
COPY nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /etc/greenstar && echo "${VERSION}" > /etc/greenstar/version
COPY --from=builder /app/dist /app/
