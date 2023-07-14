FROM node:latest AS builder

WORKDIR /app

COPY . .
RUN npm install
RUN npm run build

FROM node:latest AS runner
LABEL org.opencontainers.image.source="https://github.com/do3-2023/thomas-kube"

WORKDIR /app

ENV NODE_ENV production
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000

ENV PORT 3000
CMD ["npm", "start"]