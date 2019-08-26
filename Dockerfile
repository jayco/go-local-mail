FROM golang:1.12 AS build
WORKDIR /usr/app
COPY . /usr/app
RUN make buildlinux

FROM scratch
EXPOSE 8080
COPY --from=build /usr/app/bin/linux/local-mail /go/bin/local-mail
ENTRYPOINT [ "/go/bin/local-mail" ]
