package com.digitalvenue.http;

import io.javalin.Javalin;
import lombok.Builder;
import lombok.Data;

public class Server {

  @Data
  @Builder
  public static class Config {

    private int port;
  }

  private final Javalin app;
  private final Config config;

  public Server(final Config config) {
    this.config = config;
    app = Javalin.create();
  }

  public void start() {
    app.start(config.getPort());
  }
}
