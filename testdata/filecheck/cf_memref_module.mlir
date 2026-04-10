// CHECK: module {
// CHECK: func.func @branchy(%arg0: i1, %arg1: index) -> i32 {
// CHECK: %[[C0:.*]] = arith.constant 0 : index
// CHECK: %[[C7:.*]] = arith.constant 7 : i32
// CHECK: %[[M:.*]] = memref.alloca(%arg1) : memref<?xi32>
// CHECK: memref.store %[[C7]], %[[M]][%[[C0]]] : memref<?xi32>
// CHECK: %[[V:.*]] = memref.load %[[M]][%[[C0]]] : memref<?xi32>
// CHECK: cf.cond_br %arg0, ^bb1, ^bb2
// CHECK: ^bb1:
// CHECK: return %[[V]] : i32
// CHECK: ^bb2:
// CHECK: return %[[C7]] : i32
// CHECK: }
