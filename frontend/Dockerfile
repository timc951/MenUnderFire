# Build stage — MenUnderFire React frontend
FROM node:24-alpine AS build
WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build arguments for environment variables
ARG VITE_API_URL
ARG VITE_AUTH_PUBLISHABLE_KEY
ARG VITE_HIT_TOKEN

# Set environment variables for build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_AUTH_PUBLISHABLE_KEY=$VITE_AUTH_PUBLISHABLE_KEY
ENV VITE_HIT_TOKEN=$VITE_HIT_TOKEN

# Build the application
RUN npm run build

# Production stage - using nginx
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN chown -R nginx:nginx /var/cache/nginx /var/log/nginx /etc/nginx/conf.d && \
    touch /var/run/nginx.pid && chown nginx:nginx /var/run/nginx.pid

EXPOSE 8080

USER nginx

CMD ["nginx", "-g", "daemon off;"]
