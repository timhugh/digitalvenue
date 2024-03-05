cmake_minimum_required(VERSION 3.5)
project(hello LANGUAGES CXX)

find_package(aws-lambda-runtime REQUIRED)

add_executable(${PROJECT_NAME} "main.cpp")
target_compile_features(${PROJECT_NAME} PRIVATE cxx_std_17)
target_link_libraries(${PROJECT_NAME} PRIVATE AWS::aws-lambda-runtime)
aws_lambda_package_target(${PROJECT_NAME} NO_LIBC)

