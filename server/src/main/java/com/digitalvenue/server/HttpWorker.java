package com.digitalvenue.server;

import com.digitalvenue.common.events.Bus;
import com.digitalvenue.common.workers.Worker;
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
  private final Bus events;
  private final Javalin javalin;

  public HttpWorker(Bus events, Config config) {
    this.config = config;
    this.events = events;
    this.javalin = Javalin.create();
  }

  public void start() throws FatalException {
    javalin.start(config.getPort());
    javalin.get("/health", ctx -> ctx.status(200).result("OK"));
  }
}
