# Build stage
FROM golang AS build-env
ADD . /src
ENV CGO_ENABLED=0
RUN cd /src && go build -o /app

# Production stage
FROM scratch
COPY --from=build-env /app /

ENTRYPOINT ["/app"]
