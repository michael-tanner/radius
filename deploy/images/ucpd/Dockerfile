# Use distroless image which already includes ca-certificates
FROM gcr.io/distroless/static:nonroot

# Argument for target architecture
ARG TARGETARCH

# Set the working directory
WORKDIR /

# Copy the application binary for the specified architecture
COPY ./linux_${TARGETARCH:-amd64}/release/ucpd /

# Copy the manifest files for the built-in providers
COPY ./manifest/built-in-providers/ /manifest/built-in-providers/

# Set the user to non-root (65532:65532 is the default non-root user in distroless)
USER 65532:65532

# Set the entrypoint to the application binary
ENTRYPOINT ["/ucpd"]
