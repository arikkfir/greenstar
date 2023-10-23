# syntax=docker/dockerfile:1

### Build executable
FROM node:21.0.0 as builder
WORKDIR /workspace
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/workspace/.npm \
    npm set cache /workspace/.npm && \
    npm ci
COPY apply-patches.sh ./
RUN chmod +x ./apply-patches.sh && ./apply-patches.sh
COPY public ./public
COPY src ./src
COPY .eslintrc.cjs ./
COPY index.html ./
COPY tsconfig.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
COPY vite.docker.config.ts ./
RUN npm run build --config=vite.docker.config.ts

### Target layer
FROM nginx:1.25
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -fv /etc/nginx/conf.d/default.conf
COPY --from=builder /workspace/dist /workspace/
