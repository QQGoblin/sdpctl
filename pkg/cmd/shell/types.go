package shell

import "bytes"

type OutPut struct {
	Title    string
	NodeName string
	StdOut   *bytes.Buffer
	StdErr   *bytes.Buffer
}
