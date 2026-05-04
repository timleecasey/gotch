# Code generator

`gen/gen.ml` is an OCaml program that reads a PyTorch `Declarations.yaml`
file and emits the Go bindings under `ts/` and `libtch/`. The vendored
declarations files at `gen/pytorch/Declarations-vX.Y.Z.yaml` cover the
LibTorch versions this fork targets; `gen.ml` is wired to read the v2.10
file by default.

## Regenerating after editing `gen.ml`

You only need OCaml + dune (no PyTorch source build):

```bash
opam install dune base stdio yaml          # one-time, on a fresh machine
eval $(opam env)
dune exec gen/gen.exe                       # writes ts/*-generated.go and libtch/*-generated.*
make build && make test                     # confirm the regenerated code still compiles + tests
```

The regen is deterministic: running it on an unmodified `gen.ml` against
the vendored declarations should leave every generated file
byte-identical (`git status` clean). That property is the contract — if
you change the generator, the diff in the regenerated files must
correspond exactly to the change you made in `gen.ml`.

Files produced by the generator:

- `ts/tensor-generated.go`
- `ts/must-tensor-generated.go`
- `libtch/torch_api_generated.h`
- `libtch/torch_api_generated.cpp.h`
- `libtch/c-generated.go`

Do not hand-edit these. If you need to change one of them, change
`gen/gen.ml` and regenerate.

## In-place op leak fix

The `fixed 1` and `fixed ntensors` branches in `gen.ml` originally
emitted `ts.ctensor = *ptr` for in-place methods (anything ending in
`_`, e.g. `Uniform_`, `Normal_`, `Zero_`, `Fill_`). That overwrote the
caller's tensor handle without freeing the heap-allocated
`torch::Tensor` that the C wrapper had just produced via
`out__[0] = new torch::Tensor(self->op_(...))`, leaking one wrapper
object plus one storage refcount per call.

The current emit calls `lib.AtFree(*ptr)` instead, discarding the
duplicate handle and leaving `ts.ctensor` pointing at the original
storage (which the in-place op already mutated). See the comment block
at `gen.ml` ~line 988.

## Producing a new `Declarations.yaml` (only when bumping LibTorch)

You only need this if you're targeting a new LibTorch version that
doesn't have a vendored `Declarations-v*.yaml` file yet:

```bash
git clone -b vX.Y.Z --recurse-submodule https://github.com/pytorch/pytorch.git
mkdir pytorch-build && cd pytorch-build
cmake -DBUILD_SHARED_LIBS:BOOL=ON -DCMAKE_BUILD_TYPE:STRING=Release \
      -DPYTHON_EXECUTABLE:PATH=`which python3` \
      -DCMAKE_INSTALL_PREFIX:PATH=../pytorch-install ../pytorch
cmake --build . --target install
# Yaml is at: pytorch-install/share/ATEN/Declarations.yaml
#         or: pytorch-build/aten/src/ATen/Declarations.yaml
```

Copy that yaml to `gen/pytorch/Declarations-vX.Y.Z.yaml` and update
`gen.ml` to read it.

References:
1. <https://github.com/pytorch/pytorch/blob/master/docs/libtorch.rst>
2. <https://discuss.pytorch.org/t/compile-libtorch-c-api-from-source/81624>
3. <https://github.com/pytorch/pytorch/issues/12562>
