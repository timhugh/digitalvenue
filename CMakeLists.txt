cmake_minimum_required(VERSION 3.17)
project(digitalvenue)

add_subdirectory(lib-core)
add_subdirectory(lambda-test)
add_subdirectory(webhook-events)

enable_testing()
