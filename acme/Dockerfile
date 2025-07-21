# syntax=docker/dockerfile:1

FROM node:24-alpine AS production
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"
COPY ./tsconfig.json ./
COPY ./src ./src/
ARG NODE_ENV
RUN npm run build
ENTRYPOINT ["npm","run","generate"]
