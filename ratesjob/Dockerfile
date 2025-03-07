# syntax=docker/dockerfile:1

FROM node:23-alpine AS production
WORKDIR /app
COPY package.json package-lock.json ./
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"
RUN npm ci
COPY ./src ./src/
COPY ./tsconfig.json ./
ARG NODE_ENV
RUN npm run build
ENTRYPOINT ["npm","run","generate"]
