cmake_minimum_required(VERSION 3.17)
project(digitalvenue)

include(FetchContent)
FetchContent_Declare(json
        URL https://github.com/nlohmann/json/releases/download/v3.11.3/json.tar.xz)
FetchContent_MakeAvailable(json)

add_library(lib-core
        main.cpp
)
target_link_libraries(lib-core PUBLIC
        nlohmann_json::nlohmann_json
)
