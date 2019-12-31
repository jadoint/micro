FROM node:alpine

WORKDIR /home/app

CMD ["yarn", "start"]
EXPOSE 3000