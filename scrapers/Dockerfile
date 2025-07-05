# syntax=docker/dockerfile:1

FROM mcr.microsoft.com/playwright:v1.53.2-noble
RUN apt-get update && \
    apt-get install -y postgresql-client && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
ARG VERSION=0.0.0-notset
RUN npm pkg set version="${VERSION}"
COPY ./playwright.config.ts ./
COPY ./tsconfig.json ./
COPY ./src ./src/
ENV NODE_ENV=development
ENTRYPOINT ["npm","run","scrape"]
