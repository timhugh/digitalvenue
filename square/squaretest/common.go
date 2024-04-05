package squaretest

import "os"

var OrderRawJson, _ = os.ReadFile("squaretest/test-order-response.json")
var OrderJson = string(OrderRawJson)

var CustomerRawJson, _ = os.ReadFile("squaretest/test-customer-response.json")
var CustomerJson = string(CustomerRawJson)
