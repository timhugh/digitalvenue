#include "MockWebhookEventParser.h"
#include "MockWebhookSignatureVerifier.h"
#include "WebhookEventService.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>

TEST_CASE("successful processPaymentCreatedEvent") {
    std::string payload = "{}";
    std::string signature = "signature";

    WebhookEvent event {"payment.created" };
    MockWebhookEventParser parser(event);

    MockWebhookSignatureVerifier verifier(true);

    std::string notificationUrl = "https://example.com/events";
    WebhookEventService service(parser, verifier, notificationUrl);
    auto result = service.processPaymentCreatedEvent(payload, signature);

    REQUIRE(result.success);
    REQUIRE(result.message.empty());

    REQUIRE(parser.getPayload() == payload);

    REQUIRE(verifier.getPayload() == payload);
    REQUIRE(verifier.getSignature() == signature);
    REQUIRE(verifier.getSignatureKey() == "signature_key"); // TODO
    REQUIRE(verifier.getNotificationUrl() == notificationUrl);
}
