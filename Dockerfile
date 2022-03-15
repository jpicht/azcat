FROM golang AS build

ENV CGO_ENABLED 0
RUN go install github.com/jpicht/azcat/cmd/azblob@v0.2.2

FROM alpine

COPY --from=build /go/bin/azblob /bin

RUN ln -s /bin/azblob /bin/azcat && \
    ln -s /bin/azblob /bin/azls && \
    ln -s /bin/azblob /bin/azping && \
    ln -s /bin/azblob /bin/azput && \
    ln -s /bin/azblob /bin/azrm
