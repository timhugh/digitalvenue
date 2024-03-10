ARG BUILDER=timheuett/digitalvenue-builder:dev
FROM $BUILDER as build

COPY . /app
RUN cd /app && mkdir build && cd build && \
    cmake3 .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/out && cmake3 --build .

FROM public.ecr.aws/lambda/provided:al2023 AS runtime
COPY --from=build /app/build/lambda-test /var/runtime/bootstrap
CMD [ "hello" ]
