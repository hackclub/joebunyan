FROM golang:1.16-alpine

WORKDIR /usr/src/app

COPY . .

RUN go build -o app .

EXPOSE 3000

CMD [ "./app" ]