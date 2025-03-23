package com.digitalvenue.server;

import com.digitalvenue.http.Server;

public class Application {

  public static void main(String[] args) {
    Server.Config config = Server.Config.builder().port(8080).build();
    Server server = new Server(config);
    server.start();
  }
}
