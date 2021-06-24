FROM golang:1.16.5 AS BUILD
COPY . /app
RUN cd /app && go build

FROM gcr.io/distroless/base
COPY --from=BUILD /app/api /app
CMD ["/app"]