#include "digitalvenue/eventbus.hpp"
#include "http_worker.hpp"

int main() {
  dv::common::EventBus events;
  dv::server::HttpWorker http_worker(events);

  http_worker.Start();
}
