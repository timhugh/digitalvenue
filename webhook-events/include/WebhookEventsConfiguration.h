#pragma once

#include "Configuration.h"

static std::string WEBHOOK_NOTIFICATION_URL = "WEBHOOK_NOTIFICATION_URL";
static std::string LOG_LEVEL = "LOG_LEVEL";

static spdlog::level::level_enum logLevelFromString(const std::string &logLevel) {
    if (logLevel == "debug") {
        return spdlog::level::level_enum::debug;
    } else {
        return spdlog::level::level_enum::info;
    }
}

struct WebhookEventsConfiguration {
    std::string notificationUrl;
    spdlog::level::level_enum logLevel;

    WebhookEventsConfiguration() {
        digitalvenue::core::configuration::Environment env;
        notificationUrl = env.require(WEBHOOK_NOTIFICATION_URL);
        logLevel = logLevelFromString(env.get(LOG_LEVEL).value_or("info"));
    }
};
