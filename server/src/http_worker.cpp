#include "http_worker.hpp"
#include "digitalvenue/eventbus.hpp"
#include <crow/app.h>

namespace dv {
namespace server {

HttpWorker::HttpWorker(common::EventBus &events) : events_(events) {}

void HttpWorker::Start() {
  CROW_ROUTE(app_, "/").methods("GET"_method)(
      [] { return crow::response("Hello World!"); });
  app_.port(8080).multithreaded().run();
}

void HttpWorker::Stop() { app_.stop(); }

} // namespace server
} // namespace dv
