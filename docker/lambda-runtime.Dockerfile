ARG RUNTIME=digitalvenue/runtime-base:dev
ARG BUILD=digitalvenue/lambda-build:dev

FROM $BUILD AS build

FROM $RUNTIME
ARG LAMBDA

COPY --from=build /digitalvenue/build/$LAMBDA/$LAMBDA ./lambda
ENTRYPOINT [ "./lambda" ]
