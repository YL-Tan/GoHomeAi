# Use official Node.js image
FROM node:18-alpine

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy source code
COPY . .

# Expose the frontend port
EXPOSE 3000

# Start the frontend
CMD ["npm", "run", "dev"]
