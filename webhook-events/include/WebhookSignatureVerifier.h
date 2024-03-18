#pragma once

#include <string>

class IWebhookSignatureVerifier {
public:
    virtual ~IWebhookSignatureVerifier() = default;

    virtual bool verify(
            const std::string &requestBody,
            const std::string &signature,
            const std::string &signature_key,
            const std::string &notification_url
    ) = 0;
};

class WebhookSignatureVerifier: public IWebhookSignatureVerifier {
public:
    bool verify(
            const std::string &requestBody,
            const std::string &signature,
            const std::string &signature_key,
            const std::string &notification_url
    );
};
