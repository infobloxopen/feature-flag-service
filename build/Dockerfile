ARG CMD
ARG GO_BUILD_FLAGS
ARG REPO=github.com/Infoblox-CTO/atlas-feature-flag
ARG GOPRIVATE=github.com/Infoblox-CTO
ARG SSH_PRIVATE_KEY

# Build the manager binary
FROM golang:1.13 as builder

ARG REPO
ARG CMD
ARG GO_BUILD_FLAGS
ARG GOPRIVATE
ARG SSH_PRIVATE_KEY
WORKDIR /go/src/${REPO}

RUN mkdir -p ~/.ssh && \
	umask 0077 && \
	# echo “StrictHostKeyChecking no “ > ~/.ssh/config && \
	git config --global url."git@bitbucket.org:".insteadOf https://bitbucket.org/ && \
	git config --global url."git@github.com:".insteadOf https://github.com/ && \
	ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts && \
	ssh-keyscan github.com >> ~/.ssh/known_hosts
	

# RUN echo "${SSH_PRIVATE_KEY}" > ~/.ssh/id_rsa

# download the kubebuilder test binaries (kubebuilder, kube-apiserver, etcd, kubectl)
RUN os=$(go env GOOS) && \
	arch=$(go env GOARCH) && \
	mkdir -p /go/bin/ && \
	curl -L https://go.kubebuilder.io/dl/2.3.0/${os}/${arch} | tar -xz -C /go/bin ; \
	mv /go/bin/kubebuilder_2.3.0_${os}_${arch} /go/bin/kubebuilder 

ARG KUBEBUILDER_ASSETS=/go/bin/kubebuilder/bin

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
# RUN GOPRIVATE=${GOPRIVATE} go mod vendor

# Copy the go source
COPY . .

# Build
RUN mkdir -p /go/bin/

RUN GOPRIVATE=${GOPRIVATE} CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build ${GO_BUILD_FLAGS} -a -o /go/bin/${CMD} ./cmd/${CMD}

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM gcr.io/distroless/static:nonroot

FROM alpine:latest
ARG REPO
ARG CMD

RUN apk add --no-cache \
    bash \
    net-tools


WORKDIR /bin
COPY --from=builder /go/bin/${CMD} manager
COPY pkg/pb/*.swagger.json www/swagger.json
COPY pkg/svc/build_manifest.json .

# Go programs require libc
RUN mkdir -p /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENTRYPOINT ["/bin/manager", "--gateway.swaggerFile", "www/swagger.json"]
