#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"
#include "http_worker.hpp"
#include "test_worker.hpp"
#include <atomic>
#include <csignal>
#include <iostream>
#include <memory>
#include <thread>
#include <vector>

std::atomic<bool> shutdown_requested(false);

void signal_handler(int) { shutdown_requested.store(true); }

int main() {
  std::signal(SIGINT, signal_handler);
  std::signal(SIGTERM, signal_handler);

  dv::common::EventBus event_bus;

  std::vector<std::unique_ptr<dv::common::Worker>> workers;
  workers.emplace_back(std::make_unique<dv::server::HttpWorker>(event_bus));
  workers.emplace_back(std::make_unique<dv::server::TestWorker>(event_bus));

  for (auto &worker : workers) {
    worker->Start();
  }

  std::cout << "Workers started, waiting for shutdown signal..." << std::endl;
  while (!shutdown_requested) {
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
  }
  std::cout << "Shutdown requested, stopping workers..." << std::endl;

  for (auto &worker : workers) {
    worker->Stop();
  }
  std::cout << "Workers stopped" << std::endl;

  return 0;
}
