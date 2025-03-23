package com.digitalvenue.http;

import com.digitalvenue.common.EventBus;
import com.digitalvenue.common.Worker;
import com.digitalvenue.http.health.Health;
import io.javalin.Javalin;
import lombok.Builder;
import lombok.Data;

public class HttpWorker implements Worker {

  @Data
  @Builder
  public static class Config {

    @Builder.Default
    private int port = 8080;
  }

  private final Config config;
  private final EventBus events;
  private final Javalin javalin;

  public HttpWorker(final EventBus events, final Config config) {
    this.config = config;
    this.events = events;
    this.javalin = Javalin.create();
  }

  public void start() throws FatalException {
    javalin.start(config.getPort());
    javalin.get("/health", new Health.Get());
  }
}
