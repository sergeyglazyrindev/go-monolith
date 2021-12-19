FROM ubuntu:20.10 as builder
ARG DEBIAN_FRONTEND=noninteractive
ENV GOROOT=/usr/local/go
ENV GOPATH=/go-monolith
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get -y update \
    && apt-get -y install wget build-essential git-core golang npm libxml2-dev protobuf-compiler libprotobuf-dev \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/sergeyglazyrindev/go-monolith
RUN mkdir /go-monolith
RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz
RUN tar -xvf go1.16.4.linux-amd64.tar.gz
RUN mv go /usr/local
COPY . .
ARG GOPATH=/go
RUN make build

FROM ubuntu:20.10 as gomonolith
ARG DEBIAN_FRONTEND=noninteractive
ENV GOROOT=/usr/local/go
ENV GOPATH=/go-monolith
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get -y update \
    && apt-get -y install wget golang npm \
    && rm -rf /var/lib/apt/lists/*
#RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz
#RUN tar -xvf go1.16.4.linux-amd64.tar.gz
#RUN mv go /usr/local
COPY --from=builder /go-monolith/go-monolith /go-monolith/go-monolith
COPY configs/sqlite.yml /go-monolith/configs/go-monolith.yml
COPY configs/real_demo.yml /go-monolith/configs/demo.yml
ENTRYPOINT ["/go-monolith/go-monolith", "admin", "serve"]