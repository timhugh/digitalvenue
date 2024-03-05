FROM digitalvenue/dv-test/build:latest as build

FROM digitalvenue/runtime:latest AS runtime
COPY --from=build /app/build/hello /var/runtime/bootstrap
CMD [ "hello" ]
