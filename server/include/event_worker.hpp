#pragma once

#include "digitalvenue/eventbus.hpp"
#include "worker.hpp"

namespace dv {
namespace server {

template <typename EventType> class EventWorker : public Worker {
private:
  common::EventBus &events_;
  const common::EventBus::SubscriptionId subscription_id_;

public:
  EventWorker(common::EventBus &events) : events_(events) {}

  void Start() override {
    subscription_id_ = events_.Subscribe(this->processEvent);
  }
  void Stop() { events_.Unsubscribe(subscription_id_); }

  virtual void processEvent(const EventType &event) = 0;
};

} // namespace server
} // namespace dv
