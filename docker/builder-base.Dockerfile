FROM public.ecr.aws/amazonlinux/amazonlinux:2023.3.20240304.0

RUN dnf -y install gcc-c++ libcurl-devel cmake3 git openssl-devel zlib-devel tar unzip

RUN curl https://github.com/awslabs/aws-lambda-cpp/archive/refs/tags/v0.2.10.tar.gz -L -o /tmp/aws-lambda-cpp.tar.gz
RUN tar -xzf /tmp/aws-lambda-cpp.tar.gz -C /tmp
RUN mv /tmp/aws-lambda-cpp-0.2.10 /tmp/aws-lambda-cpp
RUN mkdir /tmp/aws-lambda-cpp/build
RUN cmake3 -S /tmp/aws-lambda-cpp -B /tmp/aws-lambda-cpp/build -DBUILD_SHARED_LIBS=ON -DCMAKE_INSTALL_PREFIX=/package
RUN cmake3 --build /tmp/aws-lambda-cpp/build --parallel --config Release
RUN cmake3 --install /tmp/aws-lambda-cpp/build
RUN rm -rf /tmp/aws-lambda-cpp

RUN curl https://github.com/fmtlib/fmt/releases/download/10.2.1/fmt-10.2.1.zip -L -o /tmp/fmt.zip
RUN unzip /tmp/fmt.zip -d /tmp
RUN mv /tmp/fmt-10.2.1 /tmp/fmt
RUN mkdir /tmp/fmt/build
RUN cmake3 -S /tmp/fmt -B /tmp/fmt/build -DBUILD_SHARED_LIBS=ON -DCMAKE_INSTALL_PREFIX=/package
RUN cmake3 --build /tmp/fmt/build --parallel --config Release --target fmt
RUN cmake3 --install /tmp/fmt/build
RUN rm -rf /tmp/fmt

RUN curl https://github.com/gabime/spdlog/archive/refs/tags/v1.13.0.tar.gz -L -o /tmp/spdlog.tar.gz
RUN tar -xzf /tmp/spdlog.tar.gz -C /tmp
RUN mv /tmp/spdlog-1.13.0 /tmp/spdlog
RUN mkdir /tmp/spdlog/build
RUN cmake3 -S /tmp/spdlog -B /tmp/spdlog/build -DSPDLOG_FMT_EXTERNAL=ON -DBUILD_SHARED_LIBS=ON -DCMAKE_INSTALL_PREFIX=/package
RUN cmake3 --build /tmp/spdlog/build --parallel --config Release --target spdlog
RUN cmake3 --install /tmp/spdlog/build
RUN rm -rf /tmp/spdlog

RUN curl https://github.com/aws/aws-sdk-cpp/archive/refs/tags/1.11.291.tar.gz -L -o /tmp/aws-sdk-cpp.tar.gz
RUN tar -xzf /tmp/aws-sdk-cpp.tar.gz -C /tmp
RUN mv /tmp/aws-sdk-cpp-1.11.291 /tmp/aws-sdk-cpp
RUN cd /tmp/aws-sdk-cpp && sh prefetch_crt_dependency.sh
RUN mkdir /tmp/aws-sdk-cpp/build
RUN cmake3 -S /tmp/aws-sdk-cpp -B /tmp/aws-sdk-cpp/build -DBUILD_ONLY="dynamodb" -DAUTORUN_UNIT_TESTS=OFF -DBUILD_SHARED_LIBS=ON -DCMAKE_INSTALL_PREFIX=/package
RUN cmake3 --build /tmp/aws-sdk-cpp/build --config Release
RUN cmake3 --install /tmp/aws-sdk-cpp/build
RUN rm -rf /tmp/aws-sdk-cpp
