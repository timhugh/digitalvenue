#include "MockWebhookEventParser.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>

using Catch::Matchers::Message;

TEST_CASE("returns stubbed result") {
    WebhookEventContainer event{"body"};

    MockWebhookEventParser parser;
    parser.stubResult(event);

    auto result = parser.parse("payload");

    REQUIRE(result.body == event.body);
}

TEST_CASE("throws stubbed exception") {
    parse_exception exception("message");

    MockWebhookEventParser parser;
    parser.stubException(exception);

    REQUIRE_THROWS_MATCHES(parser.parse("payload"),
                           parse_exception,
                           Message("Failed to parse payload: message"));
}

TEST_CASE("stores received payload") {
    MockWebhookEventParser parser;
    parser.stubResult({});

    parser.parse("payload");

    REQUIRE(parser.getPayload() == "payload");
}