#include "WebhookEventParser.h"
#include <nlohmann/json.hpp>
#include <fmt/core.h>

using json = nlohmann::json;

namespace digitalvenue::webhook_events {
    WebhookEventContainer WebhookEventParser::parse(const std::string &payload) {
        try {
            json containerJson = json::parse(payload);
            WebhookEventContainer container = containerJson.get<WebhookEventContainer>();

            json eventJson = json::parse(container.body);
            WebhookEvent data = eventJson.get<WebhookEvent>();

            data.signature = container.headers.at(SIGNATURE_HEADER_KEY);

            container.event = data;

            return container;
        }
        catch (const std::exception &e) {
            throw parse_exception(e.what());
        }
    }

    parse_exception::parse_exception(const std::string &message)
            : message(fmt::format(PARSE_EXCEPTION_FORMAT, message)) {}

    const char *parse_exception::what() const noexcept {
        return message.c_str();
    }
}
