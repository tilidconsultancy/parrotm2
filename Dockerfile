FROM golang

WORKDIR /app
COPY . .
WORKDIR /app/cmd

RUN go get
RUN go build
EXPOSE 8081
EXPOSE 9090
CMD [ "./cmd" ]