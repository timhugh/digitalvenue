#ifndef DV_SERVER_HTTP_WORKER_HPP
#define DV_SERVER_HTTP_WORKER_HPP

#include "digitalvenue/event_bus.hpp"
#include "digitalvenue/workers.hpp"
#include <crow/app.h>
#include <optional>

namespace dv {
namespace server {

class HttpWorker : public common::Worker {
public:
  HttpWorker(common::EventBus &);

  void Start() override;
  void Stop() override;

private:
  std::optional<std::thread> thread_;
  common::EventBus &events_;
  crow::SimpleApp app_;
};

} // namespace server
} // namespace dv

#endif // DV_SERVER_HTTP_WORKER_HPP
