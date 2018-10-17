FROM golang:1.11.0-alpine3.8

RUN apk add --no-cache \
	bash \
	gawk \
	git \
	grep \
	make \
	openssl \
	shadow \
	tar

ARG USER_ID
ARG USER_NAME
RUN groupadd -g "$USER_ID" "$USER_NAME"
RUN useradd -u "$USER_ID" -g "$USER_ID" "$USER_NAME"

USER $USER_ID

ENTRYPOINT ["/bin/sh", "-c"]
