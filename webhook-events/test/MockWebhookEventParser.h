#pragma once

#include "WebhookEventParser.h"

class MockWebhookEventParser: public IWebhookEventParser {
private:
    const WebhookEvent &result;
    std::string payload;

public:
    explicit MockWebhookEventParser(const WebhookEvent &result) : result(result) {}

    std::string getPayload() const {
        return payload;
    }

    WebhookEvent parse(const std::string &payload) override {
        this->payload = payload;
        return result;
    }
};