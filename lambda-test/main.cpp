#include <aws/lambda-runtime/runtime.h>
#include <spdlog/spdlog.h>
#include "LambdaResponse.h"

using namespace aws::lambda_runtime;
using namespace digitalvenue::core::lambda;

invocation_response event_handler(invocation_request const& request)
{
    spdlog::info("Request received '{}': '{}'", request.request_id, request.payload);
    LambdaResponse response{
        .statusCode = 200,
        .body = request.payload,
    };
    spdlog::info("Response '{}'", response.to_json());
    return invocation_response::success(response.to_json(), "application/json");
}

int main()
{
    run_handler(event_handler);
    return 0;
}
