FROM golang as builder
WORKDIR /pathchange
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/pathchange .


#FROM gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/git-init as git
FROM alpine
RUN addgroup -g 1000 pathchange-group && adduser -D pathchange -u 1001 -g 1000
RUN apk update && apk add git && apk add openssh
USER pathchange
WORKDIR /home/pathchange
COPY --from=builder /bin/pathchange ./
#COPY --from=git /ko-app/git-init /usr/bin/git-init
EXPOSE 8080
ENTRYPOINT ["./pathchange"]
