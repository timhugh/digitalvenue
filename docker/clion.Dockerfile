FROM digitalvenue/builder-base:dev

RUN curl https://github.com/catchorg/Catch2/archive/refs/tags/v3.5.3.tar.gz -L -o /tmp/Catch2.tar.gz
RUN tar -xzf /tmp/Catch2.tar.gz -C /tmp
RUN mv /tmp/Catch2-3.5.3 /tmp/Catch2
RUN mkdir /tmp/Catch2/build
RUN cmake3 -S /tmp/Catch2 -B /tmp/Catch2/build -DBUILD_SHARED_LIBS=ON
RUN cmake3 --build /tmp/Catch2/build --parallel --config Release
RUN cmake3 --install /tmp/Catch2/build
RUN rm -rf /tmp/Catch2
