#pragma once

#include <optional>
#include <string>

namespace digitalvenue::core::configuration {
    class missing_configuration_exception : public std::exception {
    private:
        std::string message;

    public:
        explicit missing_configuration_exception(const std::string &message): message(message) {}

        const char *what() const noexcept override {
            return message.c_str();
        }
    };

    class Environment {
    public:
        std::optional<std::string> get(const std::string &key);

        std::string require(const std::string &key);
    };
}
