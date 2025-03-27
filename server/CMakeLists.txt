if(${CMAKE_CURRENT_SOURCE_DIR} STREQUAL ${CMAKE_CURRENT_BINARY_DIR})
  message(FATAL_ERROR "Prevented in-tree build. Please create a build directory and run \"cmake -S ${CMAKE_SOURCE_DIR} -B .\" from there")
endif()

option(BUILD_SHARED_LIBS OFF)

if(CMAKE_BUILD_TYPE STREQUAL "Debug")
  add_definitions(-DDEBUG)
endif()

cmake_minimum_required(VERSION 3.20.0)
project(digitalvenue)

include(crow)

add_library(digitalvenue_server_lib
  src/http_worker.cpp
)
target_compile_features(digitalvenue_server_lib PRIVATE cxx_std_23)
target_include_directories(digitalvenue_server_lib
  PUBLIC
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
)
target_link_libraries(digitalvenue_server_lib PUBLIC
  digitalvenue::common
  Crow::Crow
)
add_library(digitalvenue::server ALIAS digitalvenue_server_lib)

add_executable(digitalvenue_server
  app/main.cpp
)
target_compile_features(digitalvenue_server PRIVATE cxx_std_23)
target_link_libraries(digitalvenue_server PRIVATE
  digitalvenue::server
)

install(TARGETS digitalvenue_server
  RUNTIME DESTINATION bin
)

include(CTest)
if(BUILD_TESTING)
  add_subdirectory(tests)
endif()
