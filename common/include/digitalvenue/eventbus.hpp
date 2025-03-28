#ifndef DV_COMMON_EVENTBUS_H
#define DV_COMMON_EVENTBUS_H

#include <cstdint>
#include <functional>
#include <mutex>
#include <shared_mutex>
#include <typeindex>
#include <unordered_map>

namespace dv {
namespace common {

class EventBus {
public:
  using SubscriptionId = uint64_t;

private:
  using Subscriber = std::function<void(const void *)>;
  struct Subscription {
    std::type_index event_type;
    SubscriptionId subscription_id;
    Subscriber subscriber;
  };

  mutable std::shared_mutex mutex_;
  std::unordered_map<SubscriptionId, Subscription> subscriptions;
  std::unordered_map<std::type_index, std::vector<SubscriptionId>>
      subscriptions_by_event;

  SubscriptionId next_id = 0;
  SubscriptionId NextId() { return ++next_id; }

public:
  template <typename EventType>
  SubscriptionId
  Subscribe(const std::function<void(const EventType &event)> &subscriber) {
    std::unique_lock<std::shared_mutex> lock(mutex_);

    auto id = NextId();
    Subscription subscription{
        typeid(EventType), id, [subscriber](const void *event_ptr) {
          subscriber(*static_cast<const EventType *>(event_ptr));
        }};
    subscriptions.emplace(id, std::move(subscription));
    subscriptions_by_event[subscription.event_type].push_back(id);
    return id;
  }

  void Unsubscribe(const SubscriptionId &subscription_id) {
    std::unique_lock<std::shared_mutex> lock(mutex_);

    auto it = subscriptions.find(subscription_id);
    if (it == subscriptions.end()) {
      return;
    }

    const Subscription &subscription = it->second;

    auto &event_subscriptions = subscriptions_by_event[subscription.event_type];
    event_subscriptions.erase(std::remove(event_subscriptions.begin(),
                                          event_subscriptions.end(),
                                          subscription_id),
                              event_subscriptions.end());
    subscriptions.erase(it);
  }

  template <typename EventType, typename... Args>
  void Emit(Args &&...args) const {
    std::shared_lock<std::shared_mutex> lock(mutex_);

    EventType event{std::forward<Args>(args)...};

    auto it = subscriptions_by_event.find(typeid(EventType));
    if (it == subscriptions_by_event.end()) {
      return;
    }

    for (auto subscription_id : it->second) {
      const auto &subscription = subscriptions.at(subscription_id);
      subscription.subscriber(&event);
    }
  }
};

} // namespace common
} // namespace dv

#endif // DV_COMMON_EVENTBUS_H
