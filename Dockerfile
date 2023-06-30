FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/

WORKDIR /app/cmd/server
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o server .

FROM tinygo/tinygo:0.28.1 as wasm-builder

WORKDIR /app
COPY cmd/wasm/main.go ./
RUN tinygo build -o main.wasm -target wasm ./main.go
RUN cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/cmd/server/server /app/server
COPY --from=wasm-builder /app/main.wasm /app/wasm/
COPY --from=wasm-builder /app/wasm_exec.js /app/wasm/

WORKDIR /app
COPY cmd/server/templates/ ./templates/
COPY cmd/server/assets/ ./assets/

EXPOSE 8080
CMD ["./server"]
