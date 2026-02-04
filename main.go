package main

import (
    "bytes"
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
    blue   = color.New(color.FgBlue).SprintFunc()
    green  = color.New(color.FgGreen).SprintFunc()
    yellow = color.New(color.FgYellow).SprintFunc()
    red    = color.New(color.FgRed).SprintFunc()

    data    string
    headers []string
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "req",
        Short: "Minimalist HTTP client",
    }

    getCmd := getCmd()
    postCmd := postCmd()

    // Define flags once per command
    getCmd.Flags().StringVarP(&data, "data", "d", "", "JSON data for request body")
    getCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Custom headers (e.g. -H 'Authorization: Bearer token')")

    postCmd.Flags().StringVarP(&data, "data", "d", "", "JSON data for request body")
    postCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Custom headers (e.g. -H 'Authorization: Bearer token')")

    rootCmd.AddCommand(getCmd, postCmd)

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
            doRequest("GET", args[0], data)
        },
    }
}

func postCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "post [url]",
        Short: "Send POST request",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            doRequest("POST", args[0], data)
        },
    }
}

func doRequest(method, url, bodyStr string) {
    var body io.Reader
    contentType := ""

    if bodyStr != "" {
        body = strings.NewReader(bodyStr)
        contentType = "application/json"
    }

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s %v\n", red("Error:"), err)
        os.Exit(1)
    }

    // Set Content-Type if body is present
    if contentType != "" {
        req.Header.Set("Content-Type", contentType)
    }

    // Apply custom headers from --header / -H flags
    for _, h := range headers {
        parts := strings.SplitN(h, ":", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            req.Header.Set(key, value)
        } else {
            fmt.Fprintf(os.Stderr, "%s Invalid header format: %s (use 'Key: Value')\n", red("Warning:"), h)
        }
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s %v\n", red("Error:"), err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s Failed to read response: %v\n", red("Error:"), err)
        os.Exit(1)
    }

    fmt.Printf("%s %s\n", green(method), blue(url))
    fmt.Printf("%s %d %s\n", yellow("Status:"), resp.StatusCode, resp.Status)

    // Pretty-print JSON if possible
    var pretty json.RawMessage
    if json.Unmarshal(bodyBytes, &pretty) == nil {
        prettyJSON, _ := json.MarshalIndent(pretty, "", "  ")
        fmt.Println(string(prettyJSON))
    } else {
        fmt.Println(string(bodyBytes))
    }
}
