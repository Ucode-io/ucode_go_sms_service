FROM golang:1.23.2 as builder

#
RUN mkdir -p $GOPATH/src/gitlab.udevs.io/ucode/ucode_go_sms_service 
WORKDIR $GOPATH/src/gitlab.udevs.io/ucode/ucode_go_sms_service

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make build && \
    mv ./bin/ucode_go_sms_service /

FROM alpine
COPY --from=builder ucode_go_sms_service .
ENTRYPOINT ["/ucode_go_sms_service"]