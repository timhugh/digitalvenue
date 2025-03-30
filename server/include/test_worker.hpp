#ifndef DV_SERVER_TEST_WORKER_HPP
#define DV_SERVER_TEST_WORKER_HPP

#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"

namespace dv {
namespace server {

struct TestEvent {
  const std::string message;
};

// class TestWorker : public common::AsyncWorker<TestEvent> {
// public:
//   TestWorker(dv::common::EventBus &events, asio::io_context &io_context)
//       : common::AsyncWorker<TestEvent>(events, io_context) {}

//   asio::awaitable<void> ProcessEvent(const TestEvent &event) override {
//     std::cout << "TestWorker received a message: " << event.message
//               << std::endl;
//     co_return;
//   }
// };

} // namespace server
} // namespace dv

#endif // DV_SERVER_TEST_WORKER_HPP
