package com.digitalvenue.server;

import com.digitalvenue.common.EventBus;
import com.digitalvenue.common.Worker;
import com.digitalvenue.common.Worker.FatalException;
import com.digitalvenue.http.HttpWorker;
import java.util.List;

public class Application {

  private final List<Worker> workers;

  public Application() {
    EventBus events = new EventBus();
    this.workers = List.of(
      new HttpWorker(events, HttpWorker.Config.builder().port(8080).build())
    );
  }

  public void start() throws FatalException {
    for (Worker worker : workers) {
      worker.start();
    }
  }

  public static void main(String[] args) {
    new Application().start();
  }
}
