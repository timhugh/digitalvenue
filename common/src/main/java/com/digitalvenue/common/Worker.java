package com.digitalvenue.common;

public interface Worker {
  public class FatalException extends RuntimeException {

    public FatalException(String message) {
      super(message);
    }
  }

  void start() throws FatalException;
}
