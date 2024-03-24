#pragma once

#include <string>

#include "WebhookEvent.h"
#include "WebhookEventParser.h"
#include "WebhookSignatureVerifier.h"

namespace digitalvenue::webhook_events {
    struct WebhookEventServiceResult {
        bool success;
        std::string message;
    };

    class WebhookEventService {
    private:
        IWebhookEventParser &eventParser;
        IWebhookSignatureVerifier &signatureVerifier;
        const std::string &notificationUrl;

    public:
        WebhookEventService(
                IWebhookEventParser &eventParser,
                IWebhookSignatureVerifier &signatureVerifier,
                const std::string &notificationUrl
        ) : eventParser(eventParser), signatureVerifier(signatureVerifier), notificationUrl(notificationUrl) {}

        WebhookEventServiceResult processPaymentCreatedEvent(const std::string &payload);
    };
}
