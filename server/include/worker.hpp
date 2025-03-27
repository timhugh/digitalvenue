#pragma once

namespace dv {
namespace server {

class worker {
public:
  virtual void start() = 0;
  virtual void stop() = 0;
};

} // namespace server
} // namespace dv
