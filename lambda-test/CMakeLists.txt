cmake_minimum_required(VERSION 3.17)

project(digitalvenue LANGUAGES CXX)

find_package(aws-lambda-runtime REQUIRED)

add_executable(lambda-test "main.cpp")
target_compile_features(lambda-test PRIVATE cxx_std_17)
target_link_libraries(lambda-test PRIVATE AWS::aws-lambda-runtime)
aws_lambda_package_target(lambda-test NO_LIBC)
