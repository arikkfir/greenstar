# syntax=docker/dockerfile:1

FROM node:24-alpine
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"
COPY ./playwright.config.ts ./
RUN npx playwright install chrome --with-deps
COPY ./tsconfig.json ./
COPY ./src ./src/
ENV NODE_ENV=development
ENTRYPOINT ["npm","run","scrape"]
