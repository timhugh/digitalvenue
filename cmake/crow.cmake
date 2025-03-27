include(asio)

include(FetchContent)
FetchContent_Declare(
  Crow
  GIT_REPOSITORY https://github.com/CrowCpp/Crow.git
  GIT_TAG v1.2.1.2
)

set(CROW_INSTALL OFF)
set(CROW_BUILD_EXAMPLES OFF)
set(CROW_BUILD_TESTS OFF)
set(CROW_ENABLE_SSL ON)
FetchContent_MakeAvailable(Crow)
