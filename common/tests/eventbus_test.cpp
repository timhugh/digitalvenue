#include "digitalvenue/eventbus.hpp"
#include <catch2/catch_test_macros.hpp>

namespace dv {
namespace test {

TEST_CASE("dv::common::eventbus") {
  dv::common::eventbus bus;

  SECTION("emits arbitrary events to subscribers") {
    struct event_struct {};
    bool event_received = false;

    bus.subscribe<event_struct>([&event_received](const event_struct &event) {
      event_received = true;
    });

    bus.emit<event_struct>();
    REQUIRE(event_received);
  }

  SECTION("emits arbitrary event parameters") {
    struct event_struct {
      int number;
      std::string message;
    };

    bool event_received = false;
    int number_received = -1;
    std::string message_received;

    bus.subscribe<event_struct>([&event_received, &number_received,
                                 &message_received](const event_struct &event) {
      event_received = true;
      number_received = event.number;
      message_received = event.message;
    });

    bus.emit<event_struct>(42, "Hello, world!");
    REQUIRE(event_received);
    REQUIRE(number_received == 42);
    REQUIRE(message_received == "Hello, world!");
  }

  SECTION("unsubscribes subscribers") {
    struct event_struct {};

    int num_events = 0;

    auto subscriber = [&num_events](const event_struct &event) {
      num_events++;
    };

    bus.subscribe<event_struct>(subscriber);
    bus.emit<event_struct>();
    REQUIRE(num_events == 1);
    bus.emit<event_struct>();
    REQUIRE(num_events == 2);
    bus.unsubscribe<event_struct>(subscriber);
    bus.emit<event_struct>();
    REQUIRE(num_events == 2);
  }
}

} // namespace test
} // namespace dv
