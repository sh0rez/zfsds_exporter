FROM golang as builder
COPY . /app
WORKDIR /app
RUN make static

FROM alpine
COPY --from=builder /app/zfsds-exporter /usr/local/bin
ENTRYPOINT ["/usr/local/bin/zfsds-exporter"]
