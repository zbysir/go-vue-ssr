FROM node:8.11-alpine

COPY ./node_modules /node_modules
COPY ./main.js /

EXPOSE 80
ENV NODE_ENV production
CMD ["node","main.js"]
