#pragma once

#include <map>
#include <string>
#include <nlohmann/json.hpp>

struct PaymentContainer {
    std::string location_id;
    std::string order_id;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(PaymentContainer, location_id, order_id);
};

struct ObjectContainer {
    PaymentContainer payment;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(ObjectContainer, payment);
};

struct PaymentCreatedEventData {
    ObjectContainer object;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(PaymentCreatedEventData, object);
};

struct WebhookEvent {
    std::string event_id;
    PaymentCreatedEventData data;
    std::string signature;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(WebhookEvent, event_id, data);
};

struct WebhookEventContainer {
    std::string body;
    std::map<std::string, std::string> headers;
    WebhookEvent event;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(WebhookEventContainer, body, headers);
};
