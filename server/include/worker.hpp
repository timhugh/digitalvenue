#pragma once

namespace dv {
namespace server {

class Worker {
public:
  virtual void Start() = 0;
  virtual void Stop() = 0;
};

} // namespace server
} // namespace dv
