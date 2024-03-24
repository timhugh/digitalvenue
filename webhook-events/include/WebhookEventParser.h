#pragma once

#include <string>
#include "WebhookEvent.h"

static const std::string SIGNATURE_HEADER_KEY = "x-square-hmacsha256-signature";
static const std::string PARSE_EXCEPTION_FORMAT = "Failed to parse payload: {}";

class parse_exception : public std::exception {
private:
    std::string message;

public:
    explicit parse_exception(const std::string &message);

    const char *what() const noexcept override;
};

class IWebhookEventParser {
public:
    virtual ~IWebhookEventParser() = default;

    virtual WebhookEventContainer parse(const std::string &payload) = 0;
};

class WebhookEventParser : public IWebhookEventParser {
public:
    WebhookEventContainer parse(const std::string &payload) override;
};
