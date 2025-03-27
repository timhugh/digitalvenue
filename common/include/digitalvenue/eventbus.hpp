#pragma once

#include <functional>
#include <typeindex>
#include <unordered_map>

namespace dv {
namespace common {

class eventbus {
public:
  template <typename EventType>
  void
  subscribe(const std::function<void(const EventType &event)> &subscriber) {
    subscriptions[typeid(EventType)].emplace_back(
        [subscriber](const void *event_ptr) {
          subscriber(*static_cast<const EventType *>(event_ptr));
        });
  }

  template <typename EventType>
  void
  unsubscribe(const std::function<void(const EventType &event)> &subscriber) {}

  template <typename EventType, typename... Args>
  void emit(Args &&...args) const {
    EventType event{std::forward<Args>(args)...};
    auto it = subscriptions.find(typeid(EventType));
    if (it == subscriptions.end()) {
      return;
    }

    for (auto &callback : it->second) {
      callback(&event);
    }
  }

private:
  std::unordered_map<std::type_index,
                     std::vector<std::function<void(const void *)>>>
      subscriptions;
};

} // namespace common
} // namespace dv
