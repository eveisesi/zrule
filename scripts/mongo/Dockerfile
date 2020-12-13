FROM mongo:4.4

WORKDIR /usr/src/mongo

ADD ./createServiceUser.js /docker-entrypoint-initdb.d

CMD ["mongod", "--auth"]