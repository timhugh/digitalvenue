if(${CMAKE_CURRENT_SOURCE_DIR} STREQUAL ${CMAKE_CURRENT_BINARY_DIR})
  message(FATAL_ERROR "Prevented in-tree build. Please create a build directory and run \"cmake -S ${CMAKE_SOURCE_DIR} -B .\" from there")
endif()

cmake_minimum_required(VERSION 3.20.0)
project(digitalvenue VERSION 0.1.0 LANGUAGES CXX)

option(BUILD_SHARED_LIBS "Build shared libraries" OFF)

if(CMAKE_BUILD_TYPE STREQUAL "Debug")
  add_definitions(-DDEBUG)
endif()

include(CTest)
enable_testing()

add_subdirectory(common)
add_subdirectory(server)

include(CPack)
set(CPACK_GENERATOR "TGZ")
set(CPACK_PACKAGE_FILE_NAME "digitalvenue")

if(CMAKE_EXPORT_COMPILE_COMMANDS AND NOT WIN32)
  add_custom_target(create_compile_commands_symlink ALL
            COMMAND ${CMAKE_COMMAND} -E create_symlink
            "${CMAKE_BINARY_DIR}/compile_commands.json"
            "${CMAKE_SOURCE_DIR}/compile_commands.json"
            COMMENT "Creating symlink for compile_commands.json"
            VERBATIM
    )
endif()
