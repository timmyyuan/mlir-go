def Sample_AddOp : Sample_Op<"add", [Pure]> {
}

def Sample_SubOp : Sample_Op<"sub",
    [Pure]> {
}

def SampleTrait : PredOpTrait<"sample trait">;

def NotAnOp : SomeOtherThing;
