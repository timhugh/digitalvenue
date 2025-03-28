#ifndef DV_COMMON_ULID_H
#define DV_COMMON_ULID_H

#include <string>

namespace dv {
namespace common {

class ULID {
public:
  static std::string Generate();
};

} // namespace common
} // namespace dv

#endif // DV_COMMON_ULID_H
