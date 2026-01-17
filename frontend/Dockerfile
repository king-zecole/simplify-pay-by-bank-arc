# Multi-stage Dockerfile for Next.js standalone production build

# Stage 1: Install dependencies (using pnpm via Corepack)
FROM node:20-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Copy package manifests and pnpm lockfile, enable corepack and install deps
COPY package.json pnpm-lock.yaml ./

RUN corepack enable && corepack prepare pnpm@latest --activate \
    && pnpm install --no-frozen-lockfile

# Stage 2: Build the app
FROM node:20-alpine AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Ensure pnpm is available in this stage as well
RUN corepack enable && corepack prepare pnpm@latest --activate

RUN pnpm run build

# Stage 3: Runtime image
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# create non-root user
RUN addgroup -S nextjs && adduser -S nextjs -G nextjs

# Copy necessary build outputs from builder
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

USER nextjs

EXPOSE 3000

CMD ["node", "server.js"]
