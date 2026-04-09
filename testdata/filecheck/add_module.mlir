// CHECK: module {
// CHECK: func.func @add(%arg0: i32) -> i32
// CHECK: %[[R0:.*]] = arith.addi %arg0, %arg0 : i32
// CHECK: return %[[R0]] : i32
// CHECK: }

