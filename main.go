package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "req",
        Short: "Minimalist HTTP client",
    }

    rootCmd.AddCommand(getCmd())
    rootCmd.AddCommand(postCmd())

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func getCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "get [url]",
        Short: "Send GET request",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("GET %s\n", args[0])
            // real logic coming later
        },
    }
}

func postCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "post [url]",
        Short: "Send POST request",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("POST %s\n", args[0])
            // real logic coming later
        },
    }
}
