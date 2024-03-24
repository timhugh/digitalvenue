ARG BUILD=digitalvenue/builder-base:dev

FROM $BUILD AS build

COPY . /digitalvenue
RUN mkdir /digitalvenue/build
RUN cmake3 -S /digitalvenue -B /digitalvenue/build -DCMAKE_INSTALL_PREFIX=/package
RUN cmake3 --build /digitalvenue/build --parallel --config Release
RUN cmake3 --install /digitalvenue/build
