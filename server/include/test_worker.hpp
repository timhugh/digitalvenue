#ifndef DV_SERVER_TEST_WORKER_HPP
#define DV_SERVER_TEST_WORKER_HPP

#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"
#include <iostream>

namespace dv {
namespace server {

struct TestEvent {
  const std::string message;
};

class TestWorker : public common::AsyncWorker<TestEvent> {
public:
  TestWorker(dv::common::EventBus &events)
      : common::AsyncWorker<TestEvent>(events) {}

  void ProcessEvent(const TestEvent &event) override {
    std::cout << "TestWorker received a message: " << event.message
              << std::endl;
  }
};

} // namespace server
} // namespace dv

#endif // DV_SERVER_TEST_WORKER_HPP
