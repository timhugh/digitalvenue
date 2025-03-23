package com.digitalvenue.http.webhook;

import com.digitalvenue.common.square.WebhookVerifier;
import io.javalin.http.Context;
import io.javalin.http.Handler;

public class WebhookEvents {

  static final String SQUARE_HMAC_SHA256_SIGNATURE_HEADER =
    "x-square-hmacsha256-signature";

  public static class Post implements Handler {

    private final String webhookUrl;

    public Post(final String webhookUrl) {
      this.webhookUrl = webhookUrl;
    }

    @Override
    public void handle(Context ctx) {
      final WebhookVerifier verifier = new WebhookVerifier(webhookUrl);
      final String signature = ctx.header(SQUARE_HMAC_SHA256_SIGNATURE_HEADER);
      final String requestBody = ctx.body();
      final String signatureKey = "abcd1234";
      try {
        if (verifier.verify(requestBody, signature, signatureKey)) {
          ctx.status(200).result("OK");
        } else {
          ctx.status(400).result("Invalid signature");
        }
      } catch (Exception e) {
        ctx.status(500).result("Internal server error");
      }
    }
  }
}
