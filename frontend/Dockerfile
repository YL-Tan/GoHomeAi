# Build stage
FROM node:20-alpine AS builder
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Development stage
FROM node:20-alpine AS development
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
ENV NODE_ENV=development
CMD ["npm", "run", "dev"]

# Production stage
FROM node:20-alpine AS production
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/.next ./.next
COPY --from=builder /usr/src/app/public ./public
COPY --from=builder /usr/src/app/package*.json ./
COPY --from=builder /usr/src/app/next.config.mjs ./
RUN npm ci --production
ENV NODE_ENV=production
USER node
EXPOSE 3000
CMD ["npm", "start"]
