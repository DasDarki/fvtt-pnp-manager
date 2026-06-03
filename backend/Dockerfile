FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
# Prepare the upload dir so a fresh named volume inherits nonroot ownership.
RUN mkdir -p /out/data/uploads

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
COPY --from=build --chown=65532:65532 /out/data/uploads /data/uploads
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
