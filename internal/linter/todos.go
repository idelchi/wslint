package linter

// TODO(Idelchi): Add more checkers! (line too long, etc.)
// TODO(Idelchi): https://medium.com/@hochiho/writing-testable-and-flexible-code-in-golang-2b1a0e66627a
// TODO(Idelchi): Table driven tests! There's lots of repetition here.
// TODO(Idelchi): Problem: if you "Lint" or "Fix" a file several times, the status and rows will be appended to the
// previous ones, which is not what you want.
// TODO(Idelchi): Bad coupling in this package. File, FileManager, FileReplacer...
// TODO(Idelchi): Mixing blank lines checker logic (if eof && line == "") into the file logic is not ideal. Blanks
// should somehow indicate what to do.
// TODO(Idelchi): Read https://twin.sh/articles/39/go-concurrency-goroutines-worker-pools-and-throttling-made-simple
// TODO(Idelchi): WorkerPool/worker could get an interface to run (Fix/Lint) instead of concrete types.
// TODO(Idelchi): Read https://medium.com/@andrewdavisescalona/testing-in-go-some-tools-you-can-use-f3e79b398d8d
// TODO(Idelchi): Reduce and/or combine Lint and Fix functions.
// TODO(Idelchi): Use t.Cleanup() in tests.
