package com.digitalvenue.common.square;

import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class WebhookVerifier {

  private final String webhookUrl;

  public WebhookVerifier(final String webhookUrl) {
    this.webhookUrl = webhookUrl;
  }

  public boolean verify(
    final String requestBody,
    final String signature,
    final String webhookSignatureKey
  ) throws NoSuchAlgorithmException, InvalidKeyException {
    final String signatureData = webhookUrl + requestBody;

    SecretKeySpec secretKeySpec = new SecretKeySpec(
      webhookSignatureKey.getBytes(StandardCharsets.UTF_8),
      "HmacSHA256"
    );

    Mac hmac = Mac.getInstance("HmacSHA256");
    hmac.init(secretKeySpec);

    byte[] hmacBytes = hmac.doFinal(
      signatureData.getBytes(StandardCharsets.UTF_8)
    );
    String computedSignature = Base64.getEncoder().encodeToString(hmacBytes);

    return computedSignature.equals(signature);
  }
}
