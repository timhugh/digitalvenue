#pragma once

#include "digitalvenue/eventbus.hpp"
#include "worker.hpp"

namespace dv {
namespace server {

template <typename EventType> class event_worker : public worker {
private:
  common::eventbus &events;

public:
  event_worker(common::eventbus &events) : events(events) {}

  void start() override { events.subscribe(this->processEvent); }
  void stop() {}

  virtual void processEvent(const EventType &event) = 0;
};

} // namespace server
} // namespace dv
