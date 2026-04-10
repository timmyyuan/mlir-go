// CHECK: module {
// CHECK: func.func @callee() -> i32
// CHECK: func.func private @caller() -> i32 {
// CHECK: %[[CALL:.*]] = call @{{callee.*}}() : () -> i32
// CHECK: return %[[CALL]] : i32
// CHECK: }
// CHECK: func.func @{{callee.*}}() -> i32
// CHECK: }
