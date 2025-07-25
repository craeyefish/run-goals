# ==========================
# 1) DEV STAGE (target: dev)
# ==========================
FROM node:18-alpine AS dev
WORKDIR /app

# Copy dependencies separately for caching
COPY package.json package-lock.json ./
RUN rm -rf node_modules
RUN npm ci

# Copy the rest of the source
COPY . .

# Expose Angular dev server port
EXPOSE 4200

# This container will run Angular in dev mode
CMD ["npm", "run", "start"]


# ============================
# 2) BUILD STAGE (target: build)
# ============================
FROM node:18-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

# Build your Angular app in production mode
RUN npm run build -- --configuration production


# =========================================
# 3) PRODUCTION STAGE (target: prod)
# =========================================
FROM nginx:1.27.5-alpine AS prod

# copy ONLY the browser build
COPY --from=build /app/dist/strava-goal/browser/ /usr/share/nginx/html

# optional custom config
COPY default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]