# syntax=docker/dockerfile:1

### Build executable
FROM node:20.5.1 as builder
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
RUN npm run build

### Target layer
FROM nginx:1.23
COPY nginx.conf /etc/nginx/nginx.conf
RUN rm -fv /etc/nginx/conf.d/default.conf
COPY --from=builder /workspace/dist /workspace/
