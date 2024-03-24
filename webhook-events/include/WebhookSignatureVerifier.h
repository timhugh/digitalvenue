#pragma once

#include <string>

struct WebhookSignatureVerifierResult {
    const bool verified;
    const std::string expectedSignature;
};

class IWebhookSignatureVerifier {
public:
    virtual ~IWebhookSignatureVerifier() = default;

    virtual const WebhookSignatureVerifierResult verify(
            const std::string &requestBody,
            const std::string &signature,
            const std::string &signature_key,
            const std::string &notification_url
    ) = 0;
};

class WebhookSignatureVerifier: public IWebhookSignatureVerifier {
public:
    const WebhookSignatureVerifierResult verify(
            const std::string &requestBody,
            const std::string &signature,
            const std::string &signature_key,
            const std::string &notification_url
    );
};
