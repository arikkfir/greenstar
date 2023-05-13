# syntax=docker/dockerfile:1

### Build executable
FROM node:19.4 as builder
WORKDIR /workspace
COPY package.json package-lock.json ./
RUN npm ci --only=production
COPY public ./public
COPY src ./src
COPY tsconfig.json ./
RUN npm run build

### Target layer
FROM nginx:1.23
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -fv /etc/nginx/conf.d/default.conf
COPY --from=builder /workspace/build /workspace/
