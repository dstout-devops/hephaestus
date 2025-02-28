FROM mcr.microsoft.com/devcontainers/base:debian

# install mastercard network root certs
USER root
COPY certs/mcNetworkRoot.crt /usr/local/share/ca-certificates
RUN update-ca-certificates

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

USER vscode

ENV NODE_EXTRA_CA_CERTS=/usr/local/share/ca-certificates/mcNetworkRoot.crt

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1