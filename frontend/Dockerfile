# syntax=docker/dockerfile:1

FROM node:24-alpine AS dependencies
WORKDIR /app
COPY package.json package-lock.json apply-patches.sh ./
RUN npm ci
RUN apk add --no-cache patch ed && chmod +x apply-patches.sh && ./apply-patches.sh && rm -f ./apply-patches.sh
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"

FROM node:24-alpine AS development
WORKDIR /app
ARG VITE_GRAPHQL_API_URL="https://api.greenstar.test/graphql"
ARG VITE_GEOAPIFY_API_KEY
ENV NODE_ENV=development
COPY --from=dependencies /app ./
COPY eslint.config.js ./
COPY index.html ./
COPY public ./public/
COPY src ./src/
COPY tsconfig.app.json ./
COPY tsconfig.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
ENTRYPOINT ["npm","run","dev"]

FROM development AS builder
ARG VITE_GRAPHQL_API_URL="https://api.greenstar.finance/graphql"
ARG VITE_GEOAPIFY_API_KEY
ENV NODE_ENV=production
RUN npm run build

FROM nginx:1.29 AS production
COPY nginx.conf /etc/nginx/nginx.conf
RUN mkdir -p /etc/greenstar && echo "${VERSION}" > /etc/greenstar/version
COPY --from=builder /app/dist /app/
