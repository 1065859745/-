FROM scratch
COPY main /
EXPOSE 5001
ENTRYPOINT ["/main"]