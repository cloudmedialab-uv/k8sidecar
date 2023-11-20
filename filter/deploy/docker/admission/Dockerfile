FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main admission/main.go

# Etapa de producción
FROM scratch
WORKDIR /app
COPY --from=build /app/main /app

CMD ["./main"]