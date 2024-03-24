#include "MockWebhookEventParser.h"
#include "MockWebhookSignatureVerifier.h"
#include "WebhookEventService.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>

using Catch::Matchers::Contains;

TEST_CASE("successful processPaymentCreatedEvent") {
    std::string payload = "{}";
    std::string signature = "signature";
    std::string notificationUrl = "https://example.com/events";

    WebhookEventContainer event{
            .body = payload,
            .headers = {
                    {SIGNATURE_HEADER_KEY, signature}
            },
            .event = {
                    .event_id = "event_id",
                    .data = {
                            .object = {
                                    .payment = {
                                            .location_id = "location_id",
                                            .order_id = "order_id"
                                    }
                            }
                    },
                    .signature = signature
            }
    };

    MockWebhookEventParser parser;
    parser.stubResult(event);
    const WebhookSignatureVerifierResult verifierResult = {true};
    MockWebhookSignatureVerifier verifier(verifierResult);

    WebhookEventService service(parser, verifier, notificationUrl);
    auto result = service.processPaymentCreatedEvent(payload);

    REQUIRE(result.success);
    REQUIRE(result.message.empty());

    REQUIRE(parser.getPayload() == payload);

    REQUIRE(verifier.getPayload() == payload);
    REQUIRE(verifier.getSignature() == signature);
    REQUIRE(verifier.getSignatureKey() == "signature_key"); // TODO
    REQUIRE(verifier.getNotificationUrl() == notificationUrl);
}

TEST_CASE("Returns failure on parse exception") {
    parse_exception exception("failed to parse");

    MockWebhookEventParser parser;
    parser.stubException(exception);

    const WebhookSignatureVerifierResult verifierResult = {true};
    MockWebhookSignatureVerifier verifier(verifierResult);

    WebhookEventService service(parser, verifier, "https://example.com/events");
    auto result = service.processPaymentCreatedEvent("{}");

    REQUIRE(!result.success);
    REQUIRE(result.message == "Failed to parse payload: failed to parse");
}

TEST_CASE("Returns failure on invalid signature") {
    MockWebhookEventParser parser;
    parser.stubResult({.event = {.data = {.object = {.payment = {.order_id = "order_id"}}}}});
    const WebhookSignatureVerifierResult verifierResult = {false};
    MockWebhookSignatureVerifier verifier(verifierResult);

    WebhookEventService service(parser, verifier, "https://example.com/events");
    auto result = service.processPaymentCreatedEvent("{}");

    REQUIRE(!result.success);
    REQUIRE(result.message == "Invalid signature");
}
