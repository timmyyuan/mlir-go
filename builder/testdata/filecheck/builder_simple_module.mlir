// CHECK: module {
// CHECK: func.func @increment(%arg0: i32) -> i32 {
// CHECK: %[[C5:.*]] = arith.constant 5 : i32
// CHECK: %[[SUM:.*]] = arith.addi %arg0, %[[C5]] : i32
// CHECK: return %[[SUM]] : i32
// CHECK: }
// CHECK: }
