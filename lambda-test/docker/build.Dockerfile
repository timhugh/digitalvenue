FROM digitalvenue/builder:latest as build

COPY . /app
RUN cd /app && mkdir build && cd build && \
    cmake3 .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/out && cmake3 --build .
