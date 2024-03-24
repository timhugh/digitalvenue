
#include <spdlog/spdlog.h>
#include "WebhookEventService.h"

namespace digitalvenue::webhook_events {
    WebhookEventServiceResult WebhookEventService::processPaymentCreatedEvent(const std::string &payload) {
        WebhookEventContainer eventContainer;
        try {
            eventContainer = eventParser.parse(payload);
        }
        catch (const parse_exception &e) {
            spdlog::warn(e.what());
            return {false, e.what()};
        }

        auto signature_key = "signature_key"; // TODO

        auto signatureVerificationResult =
                signatureVerifier.verify(eventContainer.body, eventContainer.event.signature, signature_key,
                                         notificationUrl);
        if (!signatureVerificationResult.verified) {
            spdlog::warn("Failed to validate signature for payment ID {}",
                         eventContainer.event.data.object.payment.order_id);
            spdlog::debug("Actual signature: {} expected: {}", eventContainer.event.signature,
                          signatureVerificationResult.expectedSignature);
            return {false, "Invalid signature"};
        }

        return {true, ""};
    }
}
