FROM thiagocsg/tld-golang as build
WORKDIR /app
COPY . .
WORKDIR /app/cmd
RUN go get -tags musl
RUN CC=gcc go build --ldflags '-s -w -linkmode external -extldflags "-static"' -tags musl
FROM alpine:latest
# Set TimeZoneENV TZ=America/Sao_Paulo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk add --no-cache bash
COPY --from=build /app/cmd/cmd /usr/local/bin/cmd
COPY --from=build /app/configs/ /configs/
EXPOSE 8081
EXPOSE 9090
ENTRYPOINT [ "cmd" ]