#include "digitalvenue/ulid.hpp"
#include <chrono>
#include <random>

namespace dv {
namespace common {

static constexpr char kEncodingChars[] =
    "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
static constexpr int kEncodingCharsLength = sizeof(kEncodingChars) - 1;

std::string ULID::Generate() {
  auto now = std::chrono::system_clock::now();
  auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(
                now.time_since_epoch())
                .count();

  std::random_device rd;
  std::mt19937_64 gen(rd());
  std::uniform_int_distribution<uint64_t> dis;

  std::string result(26, '0');

  for (int i = 9; i >= 0; i--) {
    result[i] = kEncodingChars[ms % kEncodingCharsLength];
    ms /= kEncodingCharsLength;
  }

  uint64_t rand1 = dis(gen);
  uint64_t rand2 = dis(gen);

  for (int i = 10; i < 18; i++) {
    result[i] = kEncodingChars[rand1 % kEncodingCharsLength];
    rand1 /= kEncodingCharsLength;
  }
  for (int i = 18; i < 26; i++) {
    result[i] = kEncodingChars[rand2 % kEncodingCharsLength];
    rand2 /= kEncodingCharsLength;
  }

  return result;
}

} // namespace common
} // namespace dv
