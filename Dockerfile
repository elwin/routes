FROM alpine:3.7

COPY heatmap /heatmap

RUN chmod +x /heatmap

CMD ["/heatmap", "--config", "/config.yml"]