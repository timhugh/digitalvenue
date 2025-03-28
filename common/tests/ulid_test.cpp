#include "digitalvenue/ulid.hpp"
#include <catch2/catch_test_macros.hpp>

namespace dv {
namespace common {

TEST_CASE("dv::common::ulid") {
  const std::string id1 = ulid::generate();
  const std::string id2 = ulid::generate();

  REQUIRE(id1 != id2);
}

} // namespace common
} // namespace dv
