#pragma once

#include "repository/LocationRepository.h"

#include "aws/core/Aws.h"
#include "aws/dynamodb/DynamoDBClient.h"
#include "aws/client/ClientConfiguration.h"

namespace digitalvenue::core::repository {
    class IDynamoDbClient {
    public:
        IDynamoDbClient(const Aws::Client::ClientConfiguration &clientConfiguration) = 0;

        virtual ~IDynamoDbClient() = default;
    };

    class DynamoDbClient : public IDynamoDbClient {
    private:
        Aws::DynamoDB::DynamoDBClient dynamoClient;
    public:
        DynamoDbClient() = default;
    };

    class DynamoLocationRepository : public ILocationRepository {
    private:
        Aws::DynamoDB::DynamoDBClient &dynamoClient;

    public:
        DynamoLocationRepository(Aws::DynamoDB::DynamoDBClient &dynamoDbClient) : dynamoClient(dynamoDbClient) {}

        std::optional<Location> find_by_square_location_id(const std::string &square_location_id) override {
            return std::nullopt;
        }

    };
}