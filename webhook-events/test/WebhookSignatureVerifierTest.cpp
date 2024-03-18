#include "WebhookSignatureVerifier.h"

#include <catch2/catch_all.hpp>
#include <catch2/matchers/catch_matchers_all.hpp>

TEST_CASE("verifies signatures") {
    WebhookSignatureVerifier verifier;

    // TODO: replace with something that's not from v1 webhooks
    std::string payload = R"({"type":"payment.created"})";
    std::string signature = "E9PjUAz21LpgB61TDgI8zhKSjap5oEEnvlljAhz3t7Q=";
    std::string signature_key = "signature_key";
    std::string notification_url = "https://example.com/events";

    REQUIRE(verifier.verify(payload, signature, signature_key, notification_url) == true);
}

TEST_CASE("fails on bad signature") {
    WebhookSignatureVerifier verifier;

    std::string payload = R"({"type":"payment.created"})";
    std::string signature = "bad_signature";
    std::string signature_key = "signature_key";
    std::string notification_url = "https://example.com/events";

    REQUIRE(verifier.verify(payload, signature, signature_key, notification_url) == false);
}
