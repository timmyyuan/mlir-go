// CHECK: module {
// CHECK: func.func @increment(%arg0: i32 {test.name = "input"}) -> i32 {
// CHECK: %[[C:.*]] = arith.constant 5 : i32
// CHECK: %[[SUM:.*]] = arith.addi %arg0, %[[C]] : i32
// CHECK: return %[[SUM]] : i32
// CHECK: }
