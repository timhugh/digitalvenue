#include "WebhookEventParser.h"
#include <nlohmann/json.hpp>

using json = nlohmann::json;

WebhookEvent WebhookEventParser::parse(const std::string &payload) {
    json data = json::parse(payload);
    return data.get<WebhookEvent>();
}
