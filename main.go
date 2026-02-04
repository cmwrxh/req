package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"

    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

var (
    blue   = color.New(color.FgBlue).SprintFunc()
    green  = color.New(color.FgGreen).SprintFunc()
    yellow = color.New(color.FgYellow).SprintFunc()
)

func main() {
    var rootCmd = &cobra.Command{Use: "req"}
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
            url := args[0]
            resp, err := http.Get(url)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: %v\n", err)
                os.Exit(1)
            }
            defer resp.Body.Close()

            body, _ := io.ReadAll(resp.Body)

            fmt.Printf("%s %s\n", green("GET"), blue(url))
            fmt.Printf("%s %d\n", yellow("Status:"), resp.StatusCode)

            // Pretty print JSON if possible
            var pretty json.RawMessage
            if json.Unmarshal(body, &pretty) == nil {
                prettyJSON, _ := json.MarshalIndent(pretty, "", "  ")
                fmt.Println(string(prettyJSON))
            } else {
                fmt.Println(string(body))
            }
        },
    }
}

// postCmd stub for now
func postCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "post [url]",
        Short: "Send POST request",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("POST %s (coming soon)\n", args[0])
        },
    }
    return cmd
}
