// Package main provides ...
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var from string = ""
var to string = ""

func init() { // {{{
	if len(os.Args) > 1 {
		from = os.Args[1]
	}
	if len(os.Args) > 2 {
		to = os.Args[2]
	}
} // }}}
func main() {
	if from == "" || to == "" {
		fmt.Println("missing arguments!\nExample:\n\tsvnmergeinfo release-1.0/ trunk")
		os.Exit(1)
	}

	revs := strings.Fields(run("svn", "mergeinfo", "--show-revs", "eligible", from, to))
	ch := make(chan string)

	//Run svn log for each revs in paralell
	for _, rev := range revs {
		go getRev(rev, ch)
	}

	//wait for the result of each svn log and print it!
	for i := 0; i < len(revs); i++ {
		fmt.Println(<-ch)
	}
}
func getRev(rev string, ch chan string) { // {{{
	ch <- run("svn", "log", "-r", rev, from)
}                                               // }}}
func run(head string, parts ...string) string { // {{{
	var err error
	var out []byte

	head, err = exec.LookPath(head)
	if err != nil {
		fmt.Printf("LookPath Error: %s", err)
	}
	out, err = exec.Command(head, parts...).CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return string(out)

} // }}}
