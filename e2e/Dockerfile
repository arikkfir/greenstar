# syntax=docker/dockerfile:1

FROM mcr.microsoft.com/playwright:v1.53.2-noble AS production
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"
COPY ./playwright.config.ts ./
COPY ./tests ./tests/
ENTRYPOINT ["npm","run","test"]
