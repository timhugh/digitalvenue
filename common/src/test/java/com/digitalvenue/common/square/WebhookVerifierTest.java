package com.digitalvenue.common.square;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.Test;

class WebhookVerifierTest {

  final String webhookUrl = "https://example.com/webhook";
  final String requestBody = "{\"hello\":\"world\"}";
  final String requestSignature =
    "2kRE5qRU2tR+tBGlDwMEw2avJ7QM4ikPYD/PJ3bd9Og=";
  final String signatureKey = "asdf1234";

  @Test
  void testVerifySuccess() throws Exception {
    final WebhookVerifier verifier = new WebhookVerifier(webhookUrl);
    assertTrue(verifier.verify(requestBody, requestSignature, signatureKey));
  }

  @Test
  void testVerifyFailure() throws Exception {
    final WebhookVerifier verifier = new WebhookVerifier(webhookUrl);
    assertFalse(
      verifier.verify(requestBody, "invalid signature", signatureKey)
    );
  }
}
