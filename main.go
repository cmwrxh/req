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
)

var data string

func main() {
    var rootCmd = &cobra.Command{
        Use:   "req",
        Short: "Minimalist HTTP client",
    }

    get := getCmd()
    post := postCmd()

    get.Flags().StringVarP(&data, "data", "d", "", "JSON data for POST/PUT")
    post.Flags().StringVarP(&data, "data", "d", "", "JSON data for POST/PUT")

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
    headers := make(map[string]string)

    if bodyStr != "" {
        body = strings.NewReader(bodyStr)
        headers["Content-Type"] = "application/json"
    }

    req, err := http.NewRequest(method, url, body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s %v\n", red("Error:"), err)
        os.Exit(1)
    }

    for k, v := range headers {
        req.Header.Set(k, v)
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
