#include <aws/lambda-runtime/runtime.h>
#include <spdlog/spdlog.h>
#include <spdlog/cfg/env.h>
#include "WebhookEventService.h"
#include "WebhookEventsConfiguration.h"

using namespace aws::lambda_runtime;

invocation_response event_handler(
        WebhookEventService &service,
        invocation_request const &request
) {
    auto result = service.processPaymentCreatedEvent(request.payload);
    if (result.success) {
        return invocation_response::success(std::string{}, "application/json");
    } else {
        return invocation_response::failure(result.message, "application/json");
    }
}

int main() {
    const WebhookEventsConfiguration config;

    spdlog::set_level(config.logLevel);
    WebhookEventParser parser;
    WebhookSignatureVerifier verifier;
    WebhookEventService service(parser, verifier, config.notificationUrl);

    auto handler = [&service](invocation_request const &request) {
        return event_handler(service, request);
    };
    run_handler(handler);
    return 0;
}
