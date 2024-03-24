ARG RUNTIME=digitalvenue/runtime-base:dev
ARG BUILD=digitalvenue/lambda-build:dev

FROM $BUILD AS build

FROM $RUNTIME
ARG LAMBDA

COPY --from=build /package/. /usr/local/
RUN ln -s /usr/local/bin/$LAMBDA /usr/local/bin/lambda
ENTRYPOINT [ "lambda" ]
