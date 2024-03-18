cmake_minimum_required(VERSION 3.17)

project(digitalvenue LANGUAGES CXX)

set(TARGET lambda-test)

add_executable(${TARGET} "main.cpp")
target_compile_features(${TARGET} PRIVATE cxx_std_17)
target_link_libraries(${TARGET} PRIVATE lib-core)
aws_lambda_package_target(${TARGET} NO_LIBC)
