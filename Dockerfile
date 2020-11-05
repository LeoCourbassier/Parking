FROM golang:1.15.3-alpine AS build
RUN apk -u add make git
RUN go get -u golang.org/x/lint/golint

WORKDIR /src
COPY --chown=nobody:nobody . .
RUN export GO111MODULE=on
RUN make

FROM alpine AS bin
COPY --from=build /src/cmd /
EXPOSE 4000
CMD [ "/br.com.mlabs" ]