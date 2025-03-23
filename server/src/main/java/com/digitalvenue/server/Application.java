package com.digitalvenue.server;

import com.digitalvenue.common.events.Bus;
import com.digitalvenue.common.workers.Worker;
import com.digitalvenue.common.workers.Worker.FatalException;
import java.util.List;

public class Application {

  private final Bus events;
  private final List<Worker> workers;

  public Application() {
    this.events = new Bus();
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
