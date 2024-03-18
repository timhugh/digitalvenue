ARG BUILD=digitalvenue/builder-base:dev

FROM $BUILD AS build

COPY . /digitalvenue
RUN mkdir /digitalvenue/build
RUN cmake3 -S /digitalvenue -B /digitalvenue/build
RUN cmake3 --build /digitalvenue/build --parallel --config Release
