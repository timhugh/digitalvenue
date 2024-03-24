#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>
#include <cstdlib>
#include "configuration/Configuration.h"

using namespace digitalvenue::core::configuration;
using Catch::Matchers::Message;

TEST_CASE("gets values from environment") {
    setenv("DV_TEST_KEY", "test_value", true);

    auto env = Environment();
    REQUIRE(env.get("DV_TEST_KEY").value() == "test_value");
    REQUIRE_FALSE(env.get("DV_NON_EXISTENT_KEY").has_value());
    REQUIRE(env.require("DV_TEST_KEY") == "test_value");
    REQUIRE_THROWS_MATCHES(
        env.require("DV_NON_EXISTENT_KEY"),
        missing_configuration_exception,
        Message("Missing required environment variable: DV_NON_EXISTENT_KEY"));
}
