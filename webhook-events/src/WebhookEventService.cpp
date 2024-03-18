
#include "WebhookEventService.h"

WebhookEventServiceResult WebhookEventService::processPaymentCreatedEvent(
        const std::string &payload,
        const std::string &signature
) {
    auto event = eventParser.parse(payload);

    auto signature_key = "signature_key"; // TODO

    if (!signatureVerifier.verify(payload, signature, signature_key, notificationUrl)) {
        return {false, "Invalid signature"};
    }

    return {true, ""};
}
