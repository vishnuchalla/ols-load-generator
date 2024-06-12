# Use Red Hat UBI minimal as the base image
FROM registry.access.redhat.com/ubi9/ubi-minimal

# Setup installation directory
WORKDIR /tmp

# Install necessary libraries and tini
RUN microdnf update -y && \
    microdnf install -y wget && \
    wget -O tini https://github.com/krallin/tini/releases/download/v0.19.0/tini && \
    chmod +x tini && \
    mv tini /usr/local/bin/tini && \
    microdnf clean all

# Copy the binary from the local file system to the image
COPY ols-load-generator /bin/ols-load-generator

# Set the display name label
LABEL io.k8s.display-name="ols-load-generator"

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/tini", "--", "sh", "-c", "ols-load-generator -D attack"]
