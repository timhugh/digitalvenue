cmake_minimum_required(VERSION 3.17)
project(digitalvenue)

option(BUILD_SHARED_LIBS "Build shared libraries" ON)

add_subdirectory(lib-core)
add_subdirectory(lambda-test)
add_subdirectory(webhook-events)

enable_testing()
