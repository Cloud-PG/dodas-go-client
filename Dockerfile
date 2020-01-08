FROM golang:1.12.1 

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN make build && export PATH=$PATH:/app
ENTRYPOINT ["/app/dodas", "--config", "/app/.dodas.yaml"]