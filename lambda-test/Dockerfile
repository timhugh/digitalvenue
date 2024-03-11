ARG BUILDER=timheuett/digitalvenue/builder:dev
ARG RUNNER=timheuett/digitalvenue/runner:dev

FROM $BUILDER as build

COPY . /app
RUN cd /app && mkdir build && cd build && \
    cmake3 .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/out && cmake3 --build .

FROM $RUNNER

COPY --from=build /app/build/lambda-test ./lambda-test
ENTRYPOINT [ "./lambda-test" ]
