package com.digitalvenue.http.health;

import io.javalin.http.Context;
import io.javalin.http.Handler;

public class Health {

  public static class Get implements Handler {

    @Override
    public void handle(Context ctx) {
      ctx.status(200).result("OK");
    }
  }
}
