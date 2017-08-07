FROM alpine:3.2

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ADD jira-time-log /jira-time-log
ADD cookies.txt /cookies.txt
EXPOSE "8080:8080"
ENTRYPOINT [ "/jira-time-log"]