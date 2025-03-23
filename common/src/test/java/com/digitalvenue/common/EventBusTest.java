package com.digitalvenue.common;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.concurrent.atomic.AtomicReference;
import org.junit.jupiter.api.Test;

public class EventBusTest {

  @Test
  void testSingleSubscriber() {
    AtomicReference<String> receivedMessage = new AtomicReference<>();
    EventBus events = new EventBus();
    events.subscribe(String.class, event -> receivedMessage.set(event));
    events.publish("Hello");
    assertEquals("Hello", receivedMessage.get());
  }

  @Test
  void testMultipleSubscribers() {
    AtomicReference<String> receivedMessage1 = new AtomicReference<>();
    AtomicReference<String> receivedMessage2 = new AtomicReference<>();
    EventBus events = new EventBus();
    events.subscribe(String.class, event -> receivedMessage1.set(event));
    events.subscribe(String.class, event -> receivedMessage2.set(event));
    events.publish("Hello");
    assertEquals("Hello", receivedMessage1.get());
    assertEquals("Hello", receivedMessage2.get());
  }

  @Test
  void testArbitraryEvents() {
    class Payload {}

    AtomicReference<Payload> receivedPayload = new AtomicReference<>();
    EventBus events = new EventBus();
    events.subscribe(Payload.class, event -> receivedPayload.set(event));

    final Payload payload = new Payload();
    events.publish(payload);
    assertEquals(payload, receivedPayload.get());
  }
}
