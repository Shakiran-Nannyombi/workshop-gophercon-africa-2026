FROM golang:1.27-alpine AS build
WORKDIR /src

# Modules come from vendor/, and the generated templ code and CSS are committed,
# so the build needs nothing from the network.
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/karibu ./cmd/web

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/karibu /app/karibu

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/karibu"]
