#pragma once

#include "digitalvenue/eventbus.hpp"
#include "worker.hpp"
#include <crow/app.h>

namespace dv {
namespace server {

class HttpWorker : public Worker {
public:
  HttpWorker(common::EventBus &);

  void Start() override;
  void Stop() override;

private:
  common::EventBus &events_;
  crow::SimpleApp app_;
};

} // namespace server
} // namespace dv
