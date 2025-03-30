#include "asio/executor_work_guard.hpp"
#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"
#include "http_worker.hpp"
#include "test_worker.hpp"
#include <asio.hpp>
#include <atomic>
#include <csignal>
#include <future>
#include <iostream>
#include <memory>
#include <thread>
#include <vector>

std::atomic<bool> shutdown_requested(false);
static constexpr int kNumThreads = 4;

void signal_handler(int signal) { shutdown_requested.store(true); }

int main() {
  std::signal(SIGINT, signal_handler);
  std::signal(SIGTERM, signal_handler);

  asio::io_context io_context;
  asio::executor_work_guard<asio::io_context::executor_type> work_guard =
      asio::make_work_guard(io_context);

  std::vector<std::thread> threads;
  for (int i = 0; i < kNumThreads; i++) {
    const int id = i;
    threads.emplace_back([&id, &io_context] {
      std::cout << "Thread  " << id << " starting" << std::endl;
      io_context.run();
      std::cout << "Thread  " << id << " finished" << std::endl;
    });
  }

  dv::common::EventBus event_bus;

  std::vector<std::unique_ptr<dv::common::Worker>> workers;
  workers.emplace_back(std::make_unique<dv::server::HttpWorker>(event_bus));
  workers.emplace_back(
      std::make_unique<dv::server::TestWorker>(event_bus, io_context));

  std::vector<std::future<void>> tasks;
  for (auto &worker : workers) {
    auto task = worker->Start();
    if (task.valid()) {
      tasks.push_back(std::move(task));
    } else {
      std::cerr << "Failed to start worker!" << std::endl;
    }
  }

  std::cout << "Workers started, waiting for shutdown signal..." << std::endl;
  while (!shutdown_requested) {
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
  }
  std::cout << "Shutdown requested, stopping workers..." << std::endl;

  for (auto &worker : workers) {
    worker->Stop();
  }
  std::cout << "Workers stopped, waiting for tasks to complete..." << std::endl;
  for (auto &task : tasks) {
    task.wait();
  }

  std::cout << "Tasks completed, stopping io context..." << std::endl;

  work_guard.reset();
  io_context.stop();
  for (auto &thread : threads) {
    thread.join();
  }

  std::cout << "All threads joined, application exiting" << std::endl;

  return 0;
}
