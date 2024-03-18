#include "WebhookEventParser.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>
#include <fstream>
#include <iostream>

TEST_CASE("parses webhook event from square docs") {
    WebhookEventParser parser;

    std::ifstream f("test/fixtures/square-webhook-event-docs-example.json");
    std::stringstream buffer;
    buffer << f.rdbuf();
    std::string payload = buffer.str();

    WebhookEvent result = parser.parse(payload);
    auto data = result.data;
    REQUIRE(result.event_id == "13b867cf-db3d-4b1c-90b6-2f32a9d78124");
    REQUIRE(data.object.payment.location_id == "S8GWD5R9QB376");
    REQUIRE(data.object.payment.order_id == "03O3USaPaAaFnI6kkwB1JxGgBsUZY");
}
