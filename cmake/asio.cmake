include(FetchContent)
FetchContent_Declare(
  asio
  GIT_REPOSITORY https://github.com/chriskohlhoff/asio.git
  GIT_TAG asio-1-30-2
)
FetchContent_GetProperties(asio)
if(NOT asio_POPULATED)
  FetchContent_Populate(asio)

  add_library(asio INTERFACE)
  set(ASIO_INCLUDE_DIR "${asio_SOURCE_DIR}/asio/include")
  target_include_directories(asio INTERFACE ${ASIO_INCLUDE_DIR})
  target_compile_definitions(asio INTERFACE
    ASIO_STANDALONE
    ASIO_NO_DEPRECATED
  )

  find_package(Threads REQUIRED)
  target_link_libraries(asio INTERFACE Threads::Threads)

  add_library(asio::asio ALIAS asio)

  # install(TARGETS asio
  #   EXPORT asio-targets
  #   INCLUDES DESTINATION include
  # )

  # install(DIRECTORY ${ASIO_INCLUDE_DIR}
  #   DESTINATION include
  #   FILES_MATCHING PATTERN "*.h" PATTERN "*.hpp"
  # )

  # install(EXPORT asio-targets
  #   FILE asio-targets.cmake
  #   NAMESPACE asio::
  #   DESTINATION lib/cmake/asio
  # )
endif()
