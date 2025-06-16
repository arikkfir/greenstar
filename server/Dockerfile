# syntax=docker/dockerfile:1

FROM node:24-alpine AS dependencies
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY ./src ./src/
COPY ./tsconfig.json ./
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"

FROM node:24-alpine AS development
WORKDIR /app
COPY --from=dependencies /app ./
ARG NODE_ENV=development
ENTRYPOINT ["npm","run","dev"]

FROM node:24-alpine AS production
WORKDIR /app
COPY --from=dependencies /app ./
ARG NODE_ENV=production
ENTRYPOINT ["npm","run","start"]
