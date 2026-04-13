// CHECK: module {
// CHECK: func.func @callee() -> i32 {
// CHECK: arith.constant 7 : i32
// CHECK: return
// CHECK: }
// CHECK: func.func @caller() -> i32 {
// CHECK: %[[CALL:.*]] = call @callee() : () -> i32
// CHECK: return %[[CALL]] : i32
// CHECK: }
// CHECK: }
