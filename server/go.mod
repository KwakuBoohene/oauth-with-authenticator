// go.mod declares the Go module: the name every Go import path in this project
// starts with, and the minimum Go version. We deliberately write NO "require"
// lines, because this project uses only the Go standard library (AGENTS.md §6).
// Why: a zero-dependency go.mod proves we aren't quietly relying on third-party code.
module oauth-with-authenticator/server

go 1.26
