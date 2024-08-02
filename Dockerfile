FROM node:16 AS ui-build

WORKDIR /app/ui

# To avoid downloading dependencies every time, copy them and download them first.
COPY ./ui/package*.json ./

RUN npm install

COPY ./ui .

ARG VITE_FIREBASE_API_KEY
ARG VITE_FIREBASE_AUTH_DOMAIN

ENV VITE_FIREBASE_API_KEY=$VITE_FIREBASE_API_KEY
ENV VITE_FIREBASE_AUTH_DOMAIN=$VITE_FIREBASE_AUTH_DOMAIN

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

WORKDIR /app

EXPOSE 8080

CMD ["./rush"]
