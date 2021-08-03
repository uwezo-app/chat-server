FROM golang:1.16-buster AS build

ENV APP_USER app
ENV APP_HOME /go/src/chat-server

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER

COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o chat-server

FROM debian:buster

ENV APP_USER app
ENV APP_HOME /go/src/chat-server

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY --chown=0:0 --from=build $APP_HOME/chat-server $APP_HOME

USER $APP_USER

EXPOSE 8080

CMD ["./chat-server"]
