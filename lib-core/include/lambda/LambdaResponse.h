#pragma once

#include "nlohmann/json.hpp"
#include <string>
#include <map>
#include <vector>

namespace digitalvenue::core::lambda {
    struct LambdaResponse {
        int statusCode;
        std::string body;
        std::map<std::string, std::string> headers = {};
        std::map<std::string, std::vector<std::string>> multiValueHeaders = {};
        bool isBase64Encoded = false;

        NLOHMANN_DEFINE_TYPE_INTRUSIVE(LambdaResponse, statusCode, body, headers, isBase64Encoded);

    public:
        [[nodiscard]] std::string to_json() const {
            nlohmann::json j = *this;
            return j.dump();
        }
    };
}
