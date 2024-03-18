FROM public.ecr.aws/amazonlinux/amazonlinux:2023.3.20240304.0

RUN dnf -y install gcc-c++ libcurl-devel cmake3 git openssl-devel zlib-devel

COPY vendor/aws-lambda-cpp /tmp/aws-lambda-cpp
RUN mkdir /tmp/aws-lambda-cpp/build
RUN cmake3 -S /tmp/aws-lambda-cpp -B /tmp/aws-lambda-cpp/build
RUN cmake3 --build /tmp/aws-lambda-cpp/build --parallel --config Release
RUN cmake3 --install /tmp/aws-lambda-cpp/build
RUN rm -rf /tmp/aws-lambda-cpp

COPY vendor/spdlog /tmp/spdlog
RUN mkdir /tmp/spdlog/build
RUN cmake3 -S /tmp/spdlog -B /tmp/spdlog/build
RUN cmake3 --build /tmp/spdlog/build --parallel --config Release
RUN cmake3 --install /tmp/spdlog/build
RUN rm -rf /tmp/spdlog

COPY vendor/aws-sdk-cpp /tmp/aws-sdk-cpp
RUN mkdir /tmp/aws-sdk-cpp/build
RUN cmake3 -S /tmp/aws-sdk-cpp -B /tmp/aws-sdk-cpp/build -DBUILD_ONLY="dynamodb" -DAUTORUN_UNIT_TESTS=OFF
RUN cmake3 --build /tmp/aws-sdk-cpp/build --config Release
RUN cmake3 --install /tmp/aws-sdk-cpp/build
RUN rm -rf /tmp/aws-sdk-cpp
