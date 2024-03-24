cmake_minimum_required(VERSION 3.17)

project(digitalvenue LANGUAGES CXX)

set(TARGET lambda-test)

add_library(${TARGET}-lib
    main.cpp
)
target_include_directories(${TARGET}-lib PUBLIC include)
target_compile_features(${TARGET}-lib PUBLIC cxx_std_17)
target_link_libraries(${TARGET}-lib PUBLIC
    lib-core
)

add_executable(${TARGET} main.cpp

)
target_compile_features(${TARGET} PUBLIC cxx_std_17)
target_link_libraries(${TARGET} PUBLIC ${TARGET}-lib)
aws_lambda_package_target(${TARGET} NO_LIBC)

if(BUILD_TESTING)
    enable_testing()
    find_package(Catch2 REQUIRED)
    add_executable(${TARGET}-test
        main.cpp
    )
    target_compile_features(${TARGET}-test PUBLIC cxx_std_17)
    target_link_libraries(${TARGET}-test PUBLIC ${TARGET}-lib Catch2::Catch2WithMain)

    add_test(NAME ${TARGET}-test COMMAND ${TARGET}-test)
endif()

install(TARGETS ${TARGET} ${TARGET}-lib
        EXPORT ${TARGET}
        LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
        ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
        RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)
