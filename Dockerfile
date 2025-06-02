FROM busybox:latest

COPY memwaste /usr/local/bin/memwaste
ENTRYPOINT ["/usr/local/bin/memwaste"]
CMD ["--amount=10M"]
