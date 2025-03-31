#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"
#include <catch2/catch_test_macros.hpp>

namespace dv::common::test {

TEST_CASE("dv::common::AsyncWorker") {
  struct TestEvent {};

  class TestWorker : public AsyncWorker<TestEvent> {
  public:
    std::vector<TestEvent> received_events;

    TestWorker(EventBus &event_bus) : AsyncWorker(event_bus) {}

  private:
    void ProcessEvent(const TestEvent &event) override {
      received_events.push_back(event);
    }
  };

  SECTION("processes events") {
    EventBus event_bus;
    TestWorker test_worker(event_bus);

    test_worker.Start();

    event_bus.Emit<TestEvent>();
    event_bus.Emit<TestEvent>();
    event_bus.Emit<TestEvent>();

    test_worker.Stop();

    REQUIRE(test_worker.received_events.size() == 3);
  }

  SECTION("does not proess events after stopping") {
    EventBus event_bus;
    TestWorker test_worker(event_bus);

    test_worker.Start();

    event_bus.Emit<TestEvent>();
    event_bus.Emit<TestEvent>();

    test_worker.Stop();

    event_bus.Emit<TestEvent>();

    REQUIRE(test_worker.received_events.size() == 2);
  }
}

} // namespace dv::common::test
