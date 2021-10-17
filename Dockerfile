FROM ubuntu:latest
WORKDIR /app
RUN apt-get -y update && apt-get install -y curl
CMD ["/app/aggregator", "--bind", "0.0.0.0:9000"]
