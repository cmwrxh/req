package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

var (
    // Colored printers
    blue   = color.New(color.FgBlue).SprintFunc()
    green  = color.New(color.FgGreen).SprintFunc()
    yellow = color.New(color.FgYellow).SprintFunc()
    red    = color.New(color.FgRed).SprintFunc()
    cyan   = color.New(color.FgCyan).SprintFunc()

    // Flags
    data    string
    headers []string
)

func main() {
    var rootCmd = &cobra.Command{
        Use:     "req",
        Short:   "Minimalist HTTP client",
        Version: "0.1.0",
    }

    // Shared flags for all subcommands
    rootCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "JSON body data")
    rootCmd.PersistentFlags().StringSliceVarP(&headers, "header", "H", []string{}, "Custom header (can be used multiple times)")

    getCmd := getCommand()
    postCmd := postCommand()

    rootCmd.AddCommand(getCmd, postCmd)

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func getCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "get [url]",
        Short: "Send a GET request",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            sendRequest("GET", args[0], "")
        },
    }
}

func postCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "post [url]",
        Short: "Send a POST request",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            sendRequest("POST", args[0], data)
        },
    }
}

func sendRequest(method, url, bodyStr string) {
    var body io.Reader
    var contentType string

    if bodyStr != "" {
        body = strings.NewReader(bodyStr)
        contentType = "application/json"
    }

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s Failed to create request: %v\n", red("Error"), err)
        os.Exit(1)
    }

    // Set Content-Type when body is present
    if contentType != "" {
        req.Header.Set("Content-Type", contentType)
    }

    // Apply custom headers
    for _, h := range headers {
        parts := strings.SplitN(h, ":", 2)
        if len(parts) != 2 {
            fmt.Fprintf(os.Stderr, "%s Invalid header format: %q (use 'Key: Value')\n", yellow("Warning"), h)
            continue
        }
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        req.Header.Set(key, value)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s Request failed: %v\n", red("Error"), err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s Failed to read response body: %v\n", red("Error"), err)
        os.Exit(1)
    }

    // Output
    fmt.Printf("%s %s\n", green(method), blue(url))

    // Colored status
    statusColor := cyan
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        statusColor = green
    } else if resp.StatusCode >= 400 {
        statusColor = red
    }
    fmt.Printf("%s %s\n", yellow("Status"), statusColor(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status)))

    // Pretty-print JSON if valid, otherwise raw
    var pretty json.RawMessage
    if json.Unmarshal(bodyBytes, &pretty) == nil {
        formatted, err := json.MarshalIndent(pretty, "", "  ")
        if err == nil {
            fmt.Println(string(formatted))
            return
        }
    }

    // Fallback: print raw response
    if len(bodyBytes) > 0 {
        fmt.Println(string(bodyBytes))
    } else {
        fmt.Println("(empty body)")
    }
}
