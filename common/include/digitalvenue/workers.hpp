#ifndef DIGITALVENUE_EVENT_WORKER_HPP
#define DIGITALVENUE_EVENT_WORKER_HPP

#include "asio/dispatch.hpp"
#include "asio/strand.hpp"
#include "asio/this_coro.hpp"
#include "asio/use_awaitable.hpp"
#include "digitalvenue/event_bus.hpp"
#include <asio.hpp>
#include <exception>
#include <future>
#include <iostream>
#include <memory>
#include <queue>
#include <stdexcept>

namespace dv {
namespace common {

class Worker {
public:
  virtual ~Worker() = default;

  virtual std::future<void> Start() = 0;
  virtual void Stop() = 0;
};

template <typename EventType> class AsyncQueue {
public:
  AsyncQueue(asio::io_context &io_context)
      : io_context_(io_context), strand_(asio::make_strand(io_context)),
        closed_(false) {}

  template <typename T> void push(T &&event) {
    asio::post(strand_, [this, event = std::forward<T>(event)]() {
      queue_.emplace(std::make_unique<EventType>(std::move(event)));

      if (wait_promise_.has_value()) {
        auto promise = std::move(wait_promise_.value());
        wait_promise_.reset();
        promise.set_value();
      }
    });
  }

  asio::awaitable<std::unique_ptr<EventType>> pop(asio::use_awaitable_t<>) {
    while (true) {
      std::unique_ptr<EventType> event;
      bool queue_closed = false;

      co_await asio::dispatch(strand_, asio::use_awaitable);

      if (!queue_.empty()) {
        event = std::move(queue_.front());
        queue_.pop();
      } else if (closed_) {
        queue_closed = true;
      } else {
        if (!wait_promise_.has_value()) {
          wait_promise_.emplace();
        }
      }

      if (event) {
        co_return std::move(event);
      }

      if (queue_closed) {
        throw std::runtime_error("queue closed");
      }

      if (!wait_promise_.has_value()) {
        auto future = wait_promise_->get_future();
        future.wait();
      }
    }
  }

  void Close() {
    asio::post(strand_, [this]() {
      closed_ = true;

      if (wait_promise_.has_value()) {
        auto promise = std::move(wait_promise_.value());
        wait_promise_.reset();
        promise.set_value();
      }
    });
  }

  bool IsClosed() const {
    std::promise<bool> close_promise;
    asio::post(strand_,
               [this, &close_promise]() { close_promise.set_value(closed_); });
    return close_promise.get_future().get();
  }

private:
  bool closed_;

  std::optional<std::promise<void>> wait_promise_;
  asio::strand<asio::io_context::executor_type> strand_;

  std::queue<std::unique_ptr<EventType>> queue_;
  asio::io_context &io_context_;
};

template <typename EventType> class AsyncWorker : public Worker {
public:
  AsyncWorker(EventBus &event_bus, asio::io_context &io_context)
      : event_bus_(event_bus), io_context_(io_context), queue_(io_context) {}

  std::future<void> Start() override {
    event_bus_.Subscribe<EventType>([this](const EventType &event) {
      std::cout << "Received event in AsyncWorker" << std::endl;
      queue_.push(event);
    });

    return asio::co_spawn(io_context_, this->ProcessEvents(), asio::use_future);
  }

  void Stop() override { queue_.Close(); }

protected:
  asio::awaitable<void> ProcessEvents() {
    std::cout << "ProcessEvents started" << std::endl;
    try {
      while (!queue_.IsClosed()) {
        std::cout << "Waiting for event..." << std::endl;
        auto event_ptr = co_await queue_.pop(asio::use_awaitable);
        std::cout << "Received event, processing..." << std::endl;

        co_await ProcessEvent(*event_ptr);
        std::cout << "Event processed" << std::endl;
      }
    } catch (const std::exception &ex) {
      std::cerr << "Encountered error while processing events: " << ex.what()
                << std::endl;
    }
    std::cout << "ProcessEvents stopped" << std::endl;
    co_return;
  }

  virtual asio::awaitable<void> ProcessEvent(const EventType &event) = 0;

private:
  EventBus &event_bus_;
  AsyncQueue<EventType> queue_;
  asio::io_context &io_context_;
};

} // namespace common
} // namespace dv

#endif // DIGITALVENUE_EVENT_WORKER_HPP
