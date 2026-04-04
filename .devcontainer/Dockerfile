FROM mcr.microsoft.com/devcontainers/go:2-1.26-bookworm

# Remove the invalid Yarn repository to fix apt update issues
RUN rm -f /etc/apt/sources.list.d/yarn.list
