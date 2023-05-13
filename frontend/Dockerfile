# syntax=docker/dockerfile:1

### Build executable
FROM node:19.4 as builder
WORKDIR /workspace
COPY package.json package-lock.json ./
RUN npm ci
COPY public ./public
COPY src ./src
COPY tsconfig.json ./
ENV PORT=80
EXPOSE 80
ENTRYPOINT ["npm"]
CMD ["run", "start"]
