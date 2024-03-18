#include <aws/lambda-runtime/runtime.h>
#include "WebhookEventService.h"

using namespace aws::lambda_runtime;

static const std::string signatureHeaderKey = "X-Square-HmacSha256-Signature";

invocation_response event_handler(
        WebhookEventService &service,
        invocation_request const &request
) {
    std::string signature; // TODO: Get signature from request.headers

    auto result = service.processPaymentCreatedEvent(request.payload, signature);
    if (result.success) {
        return invocation_response::success(std::string{}, "application/json");
    } else {
        return invocation_response::failure(result.message, "application/json");
    }
}

int main() {
    WebhookEventParser parser;
    WebhookSignatureVerifier verifier;
    WebhookEventService service(parser, verifier, "https://example.com/events");

    auto handler = [&service](invocation_request const &request) {
        return event_handler(service, request);
    };
    run_handler(handler);
    return 0;
}
