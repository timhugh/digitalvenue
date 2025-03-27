#include "digitalvenue/eventbus.hpp"
#include "http_worker.hpp"

int main() {
  dv::common::eventbus events;
  dv::server::http_worker http_worker(events);

  http_worker.start();
}
