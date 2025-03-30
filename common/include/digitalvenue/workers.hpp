#ifndef DIGITALVENUE_EVENT_WORKER_HPP
#define DIGITALVENUE_EVENT_WORKER_HPP

#include <asio.hpp>
#include <atomic>
#include <condition_variable>
#include <mutex>
#include <optional>

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

} // namespace common
} // namespace dv

#endif // DIGITALVENUE_EVENT_WORKER_HPP
