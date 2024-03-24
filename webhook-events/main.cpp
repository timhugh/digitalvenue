#include <aws/lambda-runtime/runtime.h>
#include <spdlog/spdlog.h>
#include "WebhookEventService.h"
#include "WebhookEventsConfiguration.h"
#include "LambdaResponse.h"

using namespace aws::lambda_runtime;

invocation_response event_handler(
        WebhookEventService &service,
        invocation_request const &request
) {
    auto result = service.processPaymentCreatedEvent(request.payload);
    if (result.success) {
        LambdaResponse response{.statusCode = 200};
        return invocation_response::success(
                response.to_json(),
                "application/json");
    } else {
        LambdaResponse response{.statusCode = 400, .body = result.message};
        return invocation_response::failure(
                response.to_json(),
                "application/json");
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
