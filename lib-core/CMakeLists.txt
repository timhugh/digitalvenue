cmake_minimum_required(VERSION 3.17)
project(digitalvenue)

set(TARGET lib-core)

find_package(aws-lambda-runtime REQUIRED)
find_package(fmt REQUIRED)
find_package(spdlog REQUIRED)

include(FetchContent)
FetchContent_Declare(json URL https://github.com/nlohmann/json/releases/download/v3.11.3/json.tar.xz)
FetchContent_MakeAvailable(json)

add_library(${TARGET} main.cpp
        src/configuration/Configuration.cpp
)
target_link_libraries(${TARGET} PUBLIC
        nlohmann_json::nlohmann_json
        AWS::aws-lambda-runtime
        fmt::fmt
        spdlog::spdlog
)
target_compile_features(${TARGET} PUBLIC cxx_std_17)
target_include_directories(${TARGET} PUBLIC include)

if(BUILD_TESTING)
    enable_testing()
    find_package(Catch2 REQUIRED)
    add_executable(${TARGET}-test
            test/test.cpp
            test/configuration/ConfigurationTest.cpp
    )
    target_compile_features(${TARGET}-test PUBLIC cxx_std_17)
    target_link_libraries(${TARGET}-test PUBLIC ${TARGET} Catch2::Catch2WithMain)

    add_test(NAME ${TARGET}-test COMMAND ${TARGET}-test)
endif()

install(TARGETS ${TARGET}
        EXPORT ${TARGET}
        LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
        ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
        RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)
