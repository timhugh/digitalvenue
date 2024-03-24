#include "WebhookEventParser.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>
#include <fstream>
#include <iostream>

using Catch::Matchers::StartsWith;
using Catch::Matchers::MessageMatches;

using namespace digitalvenue::webhook_events;

TEST_CASE("parses webhook event from square docs") {
    WebhookEventParser parser;

    std::ifstream f("test/fixtures/square-webhook-event-aws-proxy-payload.json");
    std::stringstream buffer;
    buffer << f.rdbuf();
    std::string payload = buffer.str();

    WebhookEventContainer result = parser.parse(payload);
    auto event = result.event;
    auto data = event.data;
    REQUIRE(event.event_id == "13b867cf-db3d-4b1c-90b6-2f32a9d78124");
    REQUIRE(data.object.payment.location_id == "S8GWD5R9QB376");
    REQUIRE(data.object.payment.order_id == "03O3USaPaAaFnI6kkwB1JxGgBsUZY");
    REQUIRE(event.signature == "abcdefg");
}

TEST_CASE("throws parse_exception on bad payload") {
    WebhookEventParser parser;

    std::string payload = "bad_payload";
    REQUIRE_THROWS_MATCHES(parser.parse(payload), parse_exception, MessageMatches(StartsWith("Failed to parse payload")));
}
