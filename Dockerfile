FROM debian:latest
LABEL maintainer="gafarov@realnoevremya.ru"
RUN apt-get update && apt-get upgrade
RUN apt-get install -y ca-certificates
RUN go get -u github.com/fogleman/rbgg
EXPOSE 80
EXPOSE 443
COPY . .
WORKDIR /build/linux
CMD [ "./main-server" ]