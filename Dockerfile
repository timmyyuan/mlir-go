FROM ubuntu:24.04

ARG LLVM_VERSION=20
ARG GO_VERSION=1.25.0

ENV MLIRGO_LLVM_VERSION=${LLVM_VERSION}
ENV DEBIAN_FRONTEND=noninteractive
ENV LLVM_CONFIG=llvm-config
ENV CC=clang-${MLIRGO_LLVM_VERSION}
ENV CXX=clang++-${MLIRGO_LLVM_VERSION}
ENV PATH=/usr/local/go/bin:/usr/local/bin:/usr/lib/llvm-${MLIRGO_LLVM_VERSION}/bin:${PATH}

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    bash \
    build-essential \
    ca-certificates \
    curl \
    git \
    gnupg \
    pkg-config \
 && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /usr/share/keyrings \
 && curl -fsSL https://apt.llvm.org/llvm-snapshot.gpg.key | gpg --dearmor -o /usr/share/keyrings/apt.llvm.org.gpg \
 && . /etc/os-release \
 && echo "deb [signed-by=/usr/share/keyrings/apt.llvm.org.gpg] http://apt.llvm.org/${VERSION_CODENAME}/ llvm-toolchain-${VERSION_CODENAME}-${LLVM_VERSION} main" > /etc/apt/sources.list.d/llvm.list \
 && apt-get update \
 && apt-get install -y --no-install-recommends \
    clang-${MLIRGO_LLVM_VERSION} \
    libc++-${MLIRGO_LLVM_VERSION}-dev \
    libc++abi-${MLIRGO_LLVM_VERSION}-dev \
    libmlir-${MLIRGO_LLVM_VERSION}-dev \
    llvm-${MLIRGO_LLVM_VERSION}-dev \
    llvm-${MLIRGO_LLVM_VERSION}-tools \
    mlir-${MLIRGO_LLVM_VERSION}-tools \
 && rm -rf /var/lib/apt/lists/* \
 && ln -sf "$(command -v llvm-config-${MLIRGO_LLVM_VERSION})" /usr/local/bin/llvm-config \
 && if command -v FileCheck-${MLIRGO_LLVM_VERSION} >/dev/null 2>&1; then ln -sf "$(command -v FileCheck-${MLIRGO_LLVM_VERSION})" /usr/local/bin/FileCheck; fi

RUN arch="$(dpkg --print-architecture)" \
 && case "$arch" in \
      amd64) go_arch="amd64" ;; \
      arm64) go_arch="arm64" ;; \
      *) echo "unsupported architecture: $arch" >&2; exit 1 ;; \
    esac \
 && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${go_arch}.tar.gz" | tar -C /usr/local -xz \
 && go version \
 && llvm-config --version \
 && FileCheck --version

WORKDIR /workspace

CMD ["bash"]
