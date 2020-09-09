

FROM golang:1.14

# Creo una carpeta para mi codigo y lo copio en ella
WORKDIR /go/src/juntosajugar
COPY . .

# Obtengo las dependencias de mi API y las Instalo
RUN go get -d -v ./...
RUN go install -v ./...

# Compilo mi imagen
RUN go build ./cmd/web

# Configuro que programa inicia el contenedor cuando arranca
ENTRYPOINT ["/go/bin/web"]
