FROM golang:1.20.3 as build

WORKDIR /app

COPY . .

RUN go mod download
RUN go test -v -parallel 500 ./...

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /server -ldflags="-s -w"

FROM gcr.io/distroless/static-debian11 as base

# # Install Doppler CLI
# RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg && \
#     curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | apt-key add - && \
#     echo "deb https://packages.doppler.com/public/cli/deb/debian any-version main" | tee /etc/apt/sources.list.d/doppler-cli.list && \
#     apt-get update && \
#     apt-get -y install doppler

COPY --from=build /server ./
 
# run binary; use vector form
CMD ["/server"]
