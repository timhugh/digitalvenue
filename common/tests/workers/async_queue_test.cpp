#include "digitalvenue/workers.hpp"
#include <catch2/catch_test_macros.hpp>
#include <optional>
#include <thread>

namespace dv::common::test {
TEST_CASE("dv::common::AsyncQueue") {
  SECTION("pushes a message to a single worker") {
    struct TestEvent {};
    AsyncQueue<const TestEvent> queue;

    std::optional<TestEvent> received_event = std::nullopt;

    auto listener = std::thread([&queue, &received_event]() {
      auto future = queue.Pop();
      future.wait();
      received_event = future.get();
    });

    queue.Push(TestEvent{});
    listener.join();

    REQUIRE(received_event.has_value());
  }

  SECTION("can handle multiple events without deadlock") {
    struct TestEvent {};
    AsyncQueue<const TestEvent> queue;

    int num_events = 0;

    auto listener = std::thread([&queue, &num_events]() {
      while (true) {
        auto future = queue.Pop();
        future.wait();
        if (auto result = future.get()) {
          num_events++;
        } else {
          break;
        }
      }
    });

    queue.Push(TestEvent{});
    queue.Push(TestEvent{});
    queue.Push(TestEvent{});
    queue.Push(TestEvent{});
    queue.Close();
    listener.join();

    REQUIRE(num_events == 4);
  }

  SECTION("late listeners still receive messages") {
    struct TestEvent {};
    AsyncQueue<const TestEvent> queue;

    queue.Push(TestEvent{});
    queue.Push(TestEvent{});
    queue.Push(TestEvent{});
    queue.Push(TestEvent{});
    queue.Close();

    int num_events = 0;

    auto listener = std::thread([&queue, &num_events]() {
      while (true) {
        auto future = queue.Pop();
        future.wait();
        if (auto result = future.get()) {
          num_events++;
        } else {
          break;
        }
      }
    });
    listener.join();

    REQUIRE(num_events == 4);
  }
}
} // namespace dv::common::test
