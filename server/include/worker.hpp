#ifndef DV_SERVER_WORKER_HPP
#define DV_SERVER_WORKER_HPP

namespace dv {
namespace server {

class Worker {
public:
  virtual void Start() = 0;
  virtual void Stop() = 0;
};

} // namespace server
} // namespace dv

#endif // DV_SERVER_WORKER_HPP
