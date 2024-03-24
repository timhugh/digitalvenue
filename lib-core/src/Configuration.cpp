#include "Configuration.h"

std::optional<std::string> digitalvenue::core::configuration::Environment::get(const std::string &key) {
    char* value = std::getenv(key.c_str());
    if (value != nullptr) {
        return {std::string(value)};
    }
    return std::nullopt;
}

std::string digitalvenue::core::configuration::Environment::require(const std::string &key) {
    auto value = get(key);
    if (value.has_value()) {
        return value.value();
    }
    throw missing_configuration_exception("Missing required environment variable: " + key);
}
