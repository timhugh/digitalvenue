#pragma once

#include <string>
#include "WebhookEvent.h"

class IWebhookEventParser {
public:
    virtual ~IWebhookEventParser() = default;

    virtual WebhookEvent parse(const std::string &payload) = 0;
};

class WebhookEventParser : public IWebhookEventParser {
public:
    WebhookEvent parse(const std::string &payload) override;
};
