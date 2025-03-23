package com.digitalvenue.common.events;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.concurrent.atomic.AtomicReference;
import org.junit.jupiter.api.Test;

public class BusTest {

  @Test
  void testSingleSubscriber() {
    AtomicReference<String> receivedMessage = new AtomicReference<>();
    Bus bus = new Bus();
    bus.subscribe(String.class, event -> receivedMessage.set(event));
    bus.publish("Hello");
    assertEquals("Hello", receivedMessage.get());
  }

  @Test
  void testMultipleSubscribers() {
    AtomicReference<String> receivedMessage1 = new AtomicReference<>();
    AtomicReference<String> receivedMessage2 = new AtomicReference<>();
    Bus bus = new Bus();
    bus.subscribe(String.class, event -> receivedMessage1.set(event));
    bus.subscribe(String.class, event -> receivedMessage2.set(event));
    bus.publish("Hello");
    assertEquals("Hello", receivedMessage1.get());
    assertEquals("Hello", receivedMessage2.get());
  }

  @Test
  void testArbitraryEvents() {
    class Payload {}

    AtomicReference<Payload> receivedPayload = new AtomicReference<>();
    Bus bus = new Bus();
    bus.subscribe(Payload.class, event -> receivedPayload.set(event));

    final Payload payload = new Payload();
    bus.publish(payload);
    assertEquals(payload, receivedPayload.get());
  }
}
