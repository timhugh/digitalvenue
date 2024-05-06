#pragma once

#include <optional>
#include <string>

namespace digitalvenue::core::repository {
    struct Location {
        const std::string location_id;
        const std::string square_location_id;
        const std::string square_signature_key;
        const std::string square_access_token;
    };

    class ILocationRepository {
    public:
        virtual ~ILocationRepository() = default;

        virtual std::optional<Location> find_by_square_location_id(const std::string &square_location_id) = 0;
    };
}
