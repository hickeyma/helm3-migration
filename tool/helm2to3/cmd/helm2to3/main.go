package main
  
import (
	"os"
)

func main() {
	cmd := newRootCmd(os.Stdout, os.Args[1:])

        if err := cmd.Execute(); err != nil {
                os.Exit(1)
        }
}
