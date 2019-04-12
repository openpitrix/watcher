# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.
FROM golang:1.11-alpine3.7 as builder

# install tools
RUN apk add --no-cache git

WORKDIR /go/src/openpitrix.io/watcher
COPY . .

RUN mkdir -p /watcher_bin
RUN CGO_ENABLED=0 GOOS=linux GOBIN=/watcher_bin go install -v -a -ldflags '-w -s' -tags netgo openpitrix.io/watcher/cmd/...

FROM alpine:3.7
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /watcher_bin/watch /usr/local/bin/

CMD ["sh"]
