#include "WebhookSignatureVerifier.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>

using namespace digitalvenue::webhook_events;

TEST_CASE("verifies signatures") {
    WebhookSignatureVerifier verifier;

    // TODO: replace with something that's not from v1 webhooks
    std::string payload = R"({"type":"payment.created"})";
    std::string signature = "E9PjUAz21LpgB61TDgI8zhKSjap5oEEnvlljAhz3t7Q=";
    std::string signature_key = "signature_key";
    std::string notification_url = "https://example.com/events";

    auto result = verifier.verify(payload, signature, signature_key, notification_url);
    REQUIRE(result.verified == true);
}

TEST_CASE("fails on bad signature") {
    WebhookSignatureVerifier verifier;

    std::string payload = R"({"type":"payment.created"})";
    std::string bad_signature = "bad_signature";
    std::string good_signature = "E9PjUAz21LpgB61TDgI8zhKSjap5oEEnvlljAhz3t7Q=";
    std::string signature_key = "signature_key";
    std::string notification_url = "https://example.com/events";

    auto result = verifier.verify(payload, bad_signature, signature_key, notification_url);
    REQUIRE(result.verified == false);
    REQUIRE(result.expectedSignature == good_signature);
}
