FROM docker:dind

RUN apk add python3 py3-pip
RUN pip install --break-system-packages awscli
