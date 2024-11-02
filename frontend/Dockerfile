# syntax=docker/dockerfile:1

FROM node:22.5.1-alpine AS dependencies
WORKDIR /app
COPY package.json package-lock.json *.sh ./
RUN npm ci
RUN sh -c '[ -f ./apply-patches.sh ] && (apk add --no-cache patch && ./apply-patches.sh) || true'

FROM node:22.5.1-alpine AS development
WORKDIR /app
ENV NODE_ENV=development
COPY --from=dependencies /app ./
COPY .eslintrc.cjs ./
COPY index.html ./
COPY public ./public/
COPY src ./src/
COPY tsconfig.app.json ./
COPY tsconfig.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
ENTRYPOINT ["npm","run","dev"]

FROM node:22.5.1-alpine AS builder
WORKDIR /app
ENV NODE_ENV=production
COPY --from=dependencies /app ./
COPY .eslintrc.cjs ./
COPY index.html ./
COPY public ./public/
COPY src ./src/
COPY tsconfig.app.json ./
COPY tsconfig.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
RUN npm run build

FROM nginx:1.25 AS production
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -fv /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /app/
