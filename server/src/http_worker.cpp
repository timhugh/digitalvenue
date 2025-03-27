#include "http_worker.hpp"
#include "digitalvenue/eventbus.hpp"
#include <crow/app.h>

namespace dv {
namespace server {

http_worker::http_worker(common::eventbus &events) : events(events) {}

void http_worker::start() {
  CROW_ROUTE(app, "/").methods("GET"_method)(
      [] { return crow::response("Hello World!"); });
  app.port(8080).multithreaded().run();
}

void http_worker::stop() { app.stop(); }

} // namespace server
} // namespace dv
