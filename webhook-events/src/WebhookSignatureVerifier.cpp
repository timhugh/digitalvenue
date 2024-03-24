#include "WebhookSignatureVerifier.h"

#include <array>
#include <iostream>
#include <openssl/sha.h>
#include <openssl/hmac.h>

#include "base64.h"

// https://stackoverflow.com/a/64570079/4467556
std::string computeSignature(const std::string &payload,
                      const std::string &signature_key) {
    std::array<unsigned char, EVP_MAX_MD_SIZE> hash;
    unsigned int hashLength;

    HMAC(EVP_sha256(), signature_key.c_str(), signature_key.size(),
         reinterpret_cast<const unsigned char *>(payload.c_str()), payload.size(),
         hash.data(), &hashLength);

    return std::string(reinterpret_cast<char *>(hash.data()), hashLength);
}

const WebhookSignatureVerifierResult WebhookSignatureVerifier::verify(const std::string &requestBody, const std::string &signature,
                                      const std::string &signature_key, const std::string &notification_url) {
    std::string payload = notification_url + requestBody;

    auto hash = computeSignature(payload, signature_key);
    auto computedSignature = base64_encode(reinterpret_cast<const unsigned char *>(hash.c_str()), hash.size());

    if (computedSignature != signature) {
        return {false, computedSignature};
    }
    return {true, {}};
}
