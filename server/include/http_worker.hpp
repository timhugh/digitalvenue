#ifndef DV_SERVER_HTTP_WORKER_HPP
#define DV_SERVER_HTTP_WORKER_HPP

#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"
#include <crow/app.h>
#include <future>

namespace dv {
namespace server {

class HttpWorker : public common::Worker {
public:
  HttpWorker(common::EventBus &);

  std::future<void> Start() override;
  void Stop() override;

private:
  common::EventBus &events_;
  crow::SimpleApp app_;
};

} // namespace server
} // namespace dv

#endif // DV_SERVER_HTTP_WORKER_HPP
