// CHECK: module {
// CHECK: func.func @select_const(%arg0: i1) -> i32 {
// CHECK: cf.cond_br %arg0, ^bb1, ^bb2
// CHECK: ^bb1:
// CHECK: arith.constant 1 : i32
// CHECK: return
// CHECK: ^bb2:
// CHECK: arith.constant 0 : i32
// CHECK: return
// CHECK: }
// CHECK: }
