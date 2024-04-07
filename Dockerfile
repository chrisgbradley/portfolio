FROM golang:1.22 AS build

COPY ./ /portfolio

WORKDIR /portfolio

RUN go mod tidy
RUN go build -o portfolio ./source

FROM alpine:3.19
RUN apk add gcompat
COPY --from=build /portfolio /portfolio
WORKDIR /portfolio
CMD ["./portfolio"]
