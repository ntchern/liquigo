FROM alpine

COPY liquigo-exec liquigo

ENTRYPOINT ["/liquigo"]
