cmake_minimum_required(VERSION 3.17)

project(digitalvenue LANGUAGES CXX)

set(TARGET "webhook-events")

find_package(OpenSSL REQUIRED)

add_library(${TARGET}-lib
        src/WebhookEventParser.cpp
        src/WebhookEventService.cpp
        src/WebhookSignatureVerifier.cpp
        include/WebhookSignatureVerifier.h
        ../lib-core/src/base64.cpp
        include/WebhookEventsConfiguration.h
)
target_include_directories(${TARGET}-lib PUBLIC include)
target_compile_features(${TARGET}-lib PUBLIC cxx_std_17)
target_link_libraries(${TARGET}-lib PUBLIC lib-core OpenSSL::Crypto)

add_executable(${TARGET} main.cpp)
target_compile_features(${TARGET} PUBLIC cxx_std_17)
target_link_libraries(${TARGET} PUBLIC ${TARGET}-lib)
aws_lambda_package_target(${TARGET} NO_LIBC)

if(BUILD_TESTING)
    enable_testing()
    find_package(Catch2 REQUIRED)
    add_executable(${TARGET}-test
            test/WebhookEventParserTest.cpp
            test/WebhookEventServiceTest.cpp
            test/MockWebhookEventParser.h
            test/WebhookSignatureVerifierTest.cpp
            test/MockWebhookSignatureVerifier.h
            test/MockWebhookEventParserTest.cpp
    )
    target_compile_features(${TARGET}-test PUBLIC cxx_std_17)
    target_link_libraries(${TARGET}-test PUBLIC ${TARGET}-lib Catch2::Catch2WithMain)

    add_test(NAME ${TARGET}-test COMMAND ${TARGET}-test)

    add_custom_target(${TARGET}-test-fixtures ALL
            COMMAND ${CMAKE_COMMAND} -E copy_directory ${CMAKE_CURRENT_SOURCE_DIR}/test/fixtures ${CMAKE_CURRENT_BINARY_DIR}/test/fixtures
            DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/test/fixtures
    )
endif()

install(TARGETS ${TARGET} ${TARGET}-lib
        EXPORT ${TARGET}
        LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
        ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
        RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)
