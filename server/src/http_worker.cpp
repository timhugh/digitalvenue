#include "http_worker.hpp"
#include "digitalvenue/event_bus.hpp"
#include "test_worker.hpp"
#include <crow/app.h>

namespace dv {
namespace server {

HttpWorker::HttpWorker(common::EventBus &events) : events_(events) {
  CROW_ROUTE(app_, "/").methods("GET"_method)(
      [] { return crow::response("Hello World!"); });
  CROW_ROUTE(app_, "/test/<int>")
      .methods("GET"_method)([&events](int num_tests) {
        for (int i = 0; i < num_tests; i++) {
          const std::string message = std::format("Test event #{}", i);
          events.Emit<TestEvent>(message);
        }
        return crow::response(std::format("Created {} test events", num_tests));
      });
  app_.port(8080).multithreaded();
}

void HttpWorker::Start() { app_.run(); }

void HttpWorker::Stop() { app_.stop(); }

} // namespace server
} // namespace dv
