#include "digitalvenue/event_bus.hpp"
#include <catch2/catch_test_macros.hpp>

namespace dv {
namespace test {

TEST_CASE("dv::common::eventbus") {
  dv::common::EventBus bus;

  SECTION("publishes emitted events to subscribers") {
    struct EventStruct {};
    bool event_received = false;

    bus.Subscribe<EventStruct>(
        [&event_received](const EventStruct &event) { event_received = true; });

    bus.Emit<EventStruct>();
    REQUIRE(event_received);
  }

  SECTION("publishes arbitrary event parameters") {
    struct EventStruct {
      int number;
      std::string message;
    };

    bool event_received = false;
    int number_received = -1;
    std::string message_received;

    bus.Subscribe<EventStruct>([&event_received, &number_received,
                                &message_received](const EventStruct &event) {
      event_received = true;
      number_received = event.number;
      message_received = event.message;
    });

    bus.Emit<EventStruct>(42, "Hello, world!");
    REQUIRE(event_received);
    REQUIRE(number_received == 42);
    REQUIRE(message_received == "Hello, world!");
  }

  SECTION("unsubscribes subscribers") {
    struct EventStruct {};

    int num_events = 0;

    auto subscriber = [&num_events](const EventStruct &event) { num_events++; };

    auto subscription = bus.Subscribe<EventStruct>(subscriber);
    bus.Emit<EventStruct>();
    REQUIRE(num_events == 1);
    bus.Emit<EventStruct>();
    REQUIRE(num_events == 2);
    bus.Unsubscribe(subscription);
    bus.Emit<EventStruct>();
    REQUIRE(num_events == 2);
  }
}

} // namespace test
} // namespace dv
