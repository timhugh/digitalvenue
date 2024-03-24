#pragma once

#include "WebhookSignatureVerifier.h"

using digitalvenue::webhook_events::IWebhookSignatureVerifier;
using digitalvenue::webhook_events::WebhookSignatureVerifierResult;

class MockWebhookSignatureVerifier : public IWebhookSignatureVerifier {
private:
    const WebhookSignatureVerifierResult result;

    std::string payload;
    std::string signature;
    std::string signature_key;
    std::string notificationUrl;

public:
    MockWebhookSignatureVerifier(const WebhookSignatureVerifierResult result) : result(result) {}

    std::string getPayload() const {
        return payload;
    }

    std::string getSignature() const {
        return signature;
    }

    std::string getSignatureKey() const {
        return signature_key;
    }

    std::string getNotificationUrl() const {
        return notificationUrl;
    }

    WebhookSignatureVerifierResult verify(
            const std::string &payload,
            const std::string &signature,
            const std::string &signature_key,
            const std::string &notificationUrl
    ) override {
        this->payload = payload;
        this->signature = signature;
        this->signature_key = signature_key;
        this->notificationUrl = notificationUrl;
        return result;
    }
};