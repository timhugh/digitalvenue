#pragma once

#include <string>

namespace digitalvenue::webhook_events {
    struct WebhookSignatureVerifierResult {
        const bool verified;
        const std::string expectedSignature;
    };

    class IWebhookSignatureVerifier {
    public:
        virtual ~IWebhookSignatureVerifier() = default;

        virtual WebhookSignatureVerifierResult verify(
                const std::string &requestBody,
                const std::string &signature,
                const std::string &signature_key,
                const std::string &notification_url
        ) = 0;
    };

    class WebhookSignatureVerifier : public IWebhookSignatureVerifier {
    public:
        WebhookSignatureVerifierResult verify(
                const std::string &requestBody,
                const std::string &signature,
                const std::string &signature_key,
                const std::string &notification_url
        ) override;
    };
}
