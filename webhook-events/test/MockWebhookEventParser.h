#pragma once

#include "WebhookEventParser.h"

#include <optional>

using namespace digitalvenue::webhook_events;

class MockWebhookEventParser: public IWebhookEventParser {
private:
    std::optional<WebhookEventContainer> result;
    std::optional<parse_exception> exception;
    std::string payload;

public:

    void stubResult(const WebhookEventContainer &result) {
        this->result = std::optional<WebhookEventContainer>(result);
    }

    void stubException(const parse_exception &exception) {
        this->exception = std::optional<parse_exception>(exception);
    }

    std::string getPayload() const {
        return payload;
    }

    WebhookEventContainer parse(const std::string &payload) override {
        this->payload = payload;

        if (exception.has_value()) {
            throw exception.value();
        }

        return result.value();
    }
};