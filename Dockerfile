FROM ysicing/god AS gobuild

LABEL maintainer="ysicing <i@ysicing.me>"

COPY . /go/src/

WORKDIR /go/src/cmd

RUN go build -o ./drone-upx

FROM ysicing/debian

RUN set -x \
    && apt-get update \
    && apt-get install --no-install-recommends --no-install-suggests -y upx \
    && rm -rf /var/lib/apt/lists/*

COPY --from=gobuild /go/src/cmd/drone-upx /bin/

RUN chmod +x /bin/drone-upx

CMD /bin/drone-upx
