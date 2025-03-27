if(${CMAKE_CURRENT_SOURCE_DIR} STREQUAL ${CMAKE_CURRENT_BINARY_DIR})
  message(FATAL_ERROR "Prevented in-tree build. Please create a build directory and run \"cmake -S ${CMAKE_SOURCE_DIR} -B .\" from there")
endif()

option(BUILD_SHARED_LIBS OFF)

if(CMAKE_BUILD_TYPE STREQUAL "Debug")
  add_definitions(-DDEBUG)
endif()

cmake_minimum_required(VERSION 3.20.0)
project(digitalvenue)

add_library(digitalvenue_common_lib
  src/keep.cpp
)
target_compile_features(digitalvenue_common_lib PRIVATE cxx_std_23)
target_include_directories(digitalvenue_common_lib
  PUBLIC
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
    $<INSTALL_INTERFACE:include>
)
target_link_libraries(digitalvenue_common_lib PUBLIC)

add_library(digitalvenue::common ALIAS digitalvenue_common_lib)

install(TARGETS digitalvenue_common_lib
  EXPORT digitalvenue-common-targets
  ARCHIVE DESTINATION lib
  LIBRARY DESTINATION lib
  RUNTIME DESTINATION bin
  INCLUDES DESTINATION include
)

install(DIRECTORY include/
  DESTINATION include
  FILES_MATCHING PATTERN "*.hpp"
)

install(EXPORT digitalvenue-common-targets
  FILE digitalvenue-common-targets.cmake
  NAMESPACE digitalvenue::
  DESTINATION lib/cmake/digitalvenue-common
)

if(BUILD_TESTING)
  add_subdirectory(tests)
endif()
