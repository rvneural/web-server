FROM debian:latest
LABEL maintainer="gafarov@realnoevremya.ru"
RUN apt-get update && apt-get upgrade
EXPOSE 80
EXPOSE 443
COPY . .
WORKDIR /build/linux
CMD [ "./main-server" ]