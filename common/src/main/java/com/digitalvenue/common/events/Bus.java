package com.digitalvenue.common.events;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

public class Bus {

  private final Map<Class<?>, List<Consumer<?>>> subscribers;

  public Bus() {
    this.subscribers = new ConcurrentHashMap<>();
  }

  @SuppressWarnings("unchecked")
  public <T> void publish(T event) {
    for (Consumer<?> subscriber : subscribers.get(event.getClass())) {
      ((Consumer<T>) subscriber).accept(event);
    }
  }

  public <T> void subscribe(Class<T> eventType, Consumer<T> subscriber) {
    subscribers
      .computeIfAbsent(eventType, k -> new ArrayList<>())
      .add(subscriber);
  }
}
