FROM alpine:3.13

# Set by docker automatically
# If building with `docker build`, make sure to set GOOS/GOARCH explicitly when calling make:
# `make compile GOOS=something GOARCH=something`
# Otherwise the makefile will not append them to the binary name and docker build will fail.
ARG TARGETOS
ARG TARGETARCH

RUN mkdir /app
WORKDIR /app

ADD --chmod=755 entrypoint.sh ./
ADD --chmod=755 bin/k8s-metadata-injection-${TARGETOS}-${TARGETARCH} ./
RUN mv k8s-metadata-injection-${TARGETOS}-${TARGETARCH} k8s-metadata-injection

CMD ["/app/entrypoint.sh"]
