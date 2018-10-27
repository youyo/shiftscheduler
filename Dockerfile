FROM golang:1.11 AS build-env
LABEL maintainer "youyo <1003ni2@gmail.com>"
ENV GO111MODULE on
ENV DIR /go/src/github.com/youyo/shiftscheduler
WORKDIR ${DIR}
ADD . ${DIR}
RUN go build -v

FROM gcr.io/distroless/base
LABEL maintainer "youyo <1003ni2@gmail.com>"
ENV PORT 1323
ENV GIN_MODE release
ENV DIR /go/src/github.com/youyo/shiftscheduler
COPY --from=build-env ${DIR}/shiftscheduler /app/shiftscheduler
EXPOSE ${PORT}/TCP
ENTRYPOINT ["/app/shiftscheduler"]
