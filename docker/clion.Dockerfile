FROM amazonlinux

RUN yum -y groupinstall "Development Tools"
RUN yum -y install \
        gcc-c++ \
        gdb \
        cmake3 \
        valgrind \
        libcurl-devel \
        ninja-build

RUN git clone https://github.com/awslabs/aws-lambda-cpp.git && \
    cd aws-lambda-cpp && mkdir build && cd build && \
    cmake3 .. -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=OFF \
    -DCMAKE_CXX_COMPILER=g++ && \
    make && make install
