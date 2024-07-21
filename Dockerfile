FROM node:16 AS ui-build

WORKDIR /app/ui

# To avoid downloading dependencies every time, copy them and download them first.
COPY ./ui/package*.json ./

RUN npm install

COPY ./ui .

RUN npm run build

FROM golang:1.22 AS server-build

WORKDIR /go/src/app

# To avoid downloading dependencies every time, copy them and download them first.
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o rush

FROM gcr.io/distroless/base-debian12

COPY --from=ui-build /app/ui/dist /app/ui/dist
COPY --from=server-build /go/src/app/rush /app
# TODO: Remove it after migrating to MongoDB.
COPY sqlite /app/sqlite

WORKDIR /app

EXPOSE 8080

CMD ["./rush"]
