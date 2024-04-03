package square

import "os"

var orderRawJson, _ = os.ReadFile("test-order-response.json")
var orderJson = string(orderRawJson)

var customerRawJson, _ = os.ReadFile("test-customer-response.json")
var customerJson = string(customerRawJson)
