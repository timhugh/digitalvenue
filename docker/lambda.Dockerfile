ARG RUNTIME=digitalvenue/runtime:dev
ARG BUILD=digitalvenue/build:dev

FROM $BUILD AS build

FROM $RUNTIME
ARG LAMBDA

COPY --from=build /digitalvenue/build/$LAMBDA/$LAMBDA ./lambda
RUN ls -la ./lambda
ENTRYPOINT [ "./lambda" ]
