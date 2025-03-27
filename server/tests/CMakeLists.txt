if(${CMAKE_CURRENT_SOURCE_DIR} STREQUAL ${CMAKE_CURRENT_BINARY_DIR})
  message(FATAL_ERROR "Prevented in-tree build. Please create a build directory and run \"cmake -S ${CMAKE_SOURCE_DIR} -B .\" from there")
endif()

cmake_minimum_required(VERSION 3.20.0)
project(digitalvenue)

include(FetchContent)
FetchContent_Declare(
    Catch2
    GIT_REPOSITORY https://github.com/catchorg/Catch2.git
    GIT_TAG v3.8.0
  )
FetchContent_MakeAvailable(Catch2)

add_executable(digitalvenue_server_tests
  test_main.cpp
)
target_link_libraries(digitalvenue_server_tests
  PRIVATE
  digitalvenue::server
  Catch2::Catch2WithMain
)

list(APPEND CMAKE_MODULE_PATH ${catch2_SOURCE_DIR}/extras)
include(Catch)
catch_discover_tests(digitalvenue_server_tests)
