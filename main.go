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
    blue   = color.New(color.FgBlue).SprintFunc()
    green  = color.New(color.FgGreen).SprintFunc()
    yellow = color.New(color.FgYellow).SprintFunc()
    red    = color.New(color.FgRed).SprintFunc()
)

var data string
var headers []string

func main() {
    var rootCmd = &cobra.Command{
        Use:   "req",
        Short: "Minimalist HTTP client",
    }

    get := getCmd()
    post := postCmd()

    get.Flags().StringVarP(&data, "data", "d", "", "JSON data for POST/PUT")
    post.Flags().StringVarP(&data, "data", "d", "", "JSON data for POST/PUT")

    get.Flags().StringSliceVarP(&headers, "header", "H", nil, "Custom header (can be repeated) -H 'Authorization: Bearer xyz'")
    post.Flags().StringSliceVarP(&headers, "header", "H", nil, "Custom header (can be repeated) -H 'Authorization: Bearer xyz'")

    rootCmd.AddCommand(get, post)

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
            doRequest("GET", args[0], "")
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

    if bodyStr != "" {
        body = strings.NewReader(bodyStr)
    }

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s %v\n", red("Error:"), err)
        os.Exit(1)
    }

    // Default Content-Type for JSON bodies
    if bodyStr != "" {
        req.Header.Set("Content-Type", "application/json")
    }

    // Apply custom headers
    for _, h := range headers {
        parts := strings.SplitN(h, ":", 2)
        if len(parts) != 2 {
            fmt.Fprintf(os.Stderr, "%s invalid header format: %s (use Key: Value)\n", red("Error:"), h)
            os.Exit(1)
        }
        key := strings.TrimSpace(parts[0])
        val := strings.TrimSpace(parts[1])
        req.Header.Set(key, val)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s %v\n", red("Error:"), err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    bodyBytes, _ := io.ReadAll(resp.Body)

    fmt.Printf("%s %s\n", green(method), blue(url))
    fmt.Printf("%s %d\n", yellow("Status:"), resp.StatusCode)

    // Pretty print JSON response if possible
    var pretty json.RawMessage
    if err := json.Unmarshal(bodyBytes, &pretty); err == nil {
        prettyJSON, _ := json.MarshalIndent(pretty, "", "  ")
        fmt.Println(string(prettyJSON))
    } else {
        fmt.Println(string(bodyBytes))
    }
}
