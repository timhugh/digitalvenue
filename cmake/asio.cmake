CPMAddPackage("gh:chriskohlhoff/asio#asio-1-30-2")
find_package(Threads REQUIRED)
if(asio_ADDED)
  add_library(asio INTERFACE)
  target_include_directories(asio SYSTEM INTERFACE ${asio_SOURCE_DIR}/asio/include)
  target_compile_definitions(asio INTERFACE ASIO_STANDALONE ASIO_NO_DEPRECATED)
  target_link_libraries(asio INTERFACE Threads::Threads)
  add_library(asio::asio ALIAS asio)
endif()
