#pragma once

#include "digitalvenue/eventbus.hpp"
#include "worker.hpp"
#include <crow/app.h>

namespace dv {
namespace server {

class http_worker : public worker {
public:
  http_worker(common::eventbus &);

  void start() override;
  void stop() override;

private:
  common::eventbus &events;
  crow::SimpleApp app;
};

} // namespace server
} // namespace dv
