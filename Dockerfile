FROM alpine:latest
MAINTAINER OpenSource PF <opensource@postfinance.ch>

RUN apk --no-cache add ca-certificates curl
COPY kubenurse /bin/kubenurse

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/bin/kubenurse"]
