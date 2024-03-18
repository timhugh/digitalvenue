cmake_minimum_required(VERSION 3.17)
project(digitalvenue)

set(TARGET lib-core)

find_package(aws-lambda-runtime REQUIRED)
find_package(spdlog REQUIRED)

include(FetchContent)
FetchContent_Declare(json URL https://github.com/nlohmann/json/releases/download/v3.11.3/json.tar.xz)
FetchContent_MakeAvailable(json)

add_library(${TARGET} main.cpp)
target_link_libraries(${TARGET} PUBLIC
        nlohmann_json::nlohmann_json
        AWS::aws-lambda-runtime
        spdlog::spdlog
)
target_include_directories(${TARGET} PUBLIC include)

if(BUILD_TESTING)
    enable_testing()
    find_package(Catch2 REQUIRED)
    add_executable(${TARGET}-test
            test/test.cpp
    )
    target_link_libraries(${TARGET}-test PRIVATE ${TARGET} Catch2::Catch2WithMain)

    add_test(NAME ${TARGET}-test COMMAND ${TARGET}-test)
endif()
