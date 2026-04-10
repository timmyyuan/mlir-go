// CHECK: module {
// CHECK: func.func @select_dim(%arg0: i1, %arg1: index) -> index {
// CHECK: %[[C0:.*]] = arith.constant 0 : index
// CHECK: %[[T:.*]] = tensor.empty(%arg1) : tensor<?xi32>
// CHECK: %[[D:.*]] = tensor.dim %[[T]], %[[C0]] : tensor<?xi32>
// CHECK: %[[R:.*]] = scf.if %arg0 -> (index) {
// CHECK: scf.yield %[[D]] : index
// CHECK: } else {
// CHECK: scf.yield %arg1 : index
// CHECK: }
// CHECK: return %[[R]] : index
// CHECK: }
