#include <crow.h>

int main() {
  crow::SimpleApp app;
  CROW_ROUTE(app, "/").methods("GET"_method)(
      [] { return crow::response("Hello World!"); });
  app.port(8080).multithreaded().run();
}
