FROM golang

WORKDIR /app
COPY . .
WORKDIR /app/cmd

RUN go get
RUN go build
EXPOSE 80
CMD [ "./cmd" ]