#ifndef DIGITALVENUE_EVENT_WORKER_HPP
#define DIGITALVENUE_EVENT_WORKER_HPP

#include "digitalvenue/event_bus.hpp"
#include <atomic>
#include <condition_variable>
#include <future>
#include <mutex>
#include <optional>
#include <thread>

namespace dv {
namespace common {

class Worker {
public:
  virtual ~Worker() = default;

  virtual void Start() = 0;
  virtual void Stop() = 0;
};

template <typename EventType> class AsyncQueue {
public:
  AsyncQueue() : closed_(false) {}

  void Push(const EventType &event) {
    {
      std::lock_guard<std::mutex> lock(mutex_);
      queue_.push(event);
    }
    condition_.notify_one();
  }

  std::future<std::optional<const EventType>> Pop() {
    return std::async(
        std::launch::async, [this]() -> std::optional<const EventType> {
          std::unique_lock<std::mutex> lock(mutex_);

          condition_.wait(lock,
                          [this]() { return !queue_.empty() || closed_; });

          if (queue_.empty() && closed_) {
            return std::nullopt;
          }

          const EventType event = queue_.front();
          queue_.pop();
          return std::make_optional(event);
        });
  }

  void Close() {
    closed_.store(true);
    condition_.notify_all();
  }

  bool IsClosed() { return closed_.load(); }

private:
  std::atomic_bool closed_;
  std::queue<EventType> queue_;
  std::mutex mutex_;
  std::condition_variable condition_;
};

template <typename EventType> class AsyncWorker : public Worker {
public:
  AsyncWorker(EventBus &event_bus) : event_bus_(event_bus) {}

  void Start() override {
    subscription_id_ = event_bus_.Subscribe<EventType>(
        [this](const EventType &event) { queue_.Push(event); });

    thread_ = std::make_optional<std::thread>([this] { ProcessEvents(); });
  }

  void Stop() override {
    event_bus_.Unsubscribe(subscription_id_);
    queue_.Close();
    if (thread_ && thread_->joinable()) {
      thread_->join();
    }
  }

protected:
  virtual void ProcessEvent(const EventType &event) = 0;

private:
  void ProcessEvents() {
    while (true) {
      auto future = queue_.Pop();
      auto opt_event = future.get();
      if (!opt_event && queue_.IsClosed()) {
        break;
      }

      ProcessEvent(*opt_event);
    }
  }

  EventBus &event_bus_;
  EventBus::SubscriptionId subscription_id_;

  AsyncQueue<EventType> queue_;
  std::optional<std::thread> thread_;
};

} // namespace common
} // namespace dv

#endif // DIGITALVENUE_EVENT_WORKER_HPP
