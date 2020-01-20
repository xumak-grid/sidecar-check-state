FROM alpine
LABEL maintainer="ehernandez@xumak.com"
WORKDIR /app
EXPOSE 9090
COPY bin/check-state ./check-state
CMD [ "/app/check-state" ]
