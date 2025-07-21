package main

import (
    "fmt"
    "net/http"
    "os"

    "dotcom/internal/generator"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run . generate")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "generate":
        cfg := generator.Config{
            ContentDir:  getenv("CONTENT_DIR", "content"),
            TemplateDir: getenv("TEMPLATE_DIR", "templates"),
            StaticDir:   getenv("STATIC_DIR", "static"),
            OutputDir:   getenv("OUTPUT_DIR", "public"),
            BaseURL:     getenv("BASE_URL", ""),
        }

        fmt.Println("generating site...")
        if err := generator.Run(cfg); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("done.")

    case "serve":
        outputDir := getenv("OUTPUT_DIR", "public")
        addr := getenv("ADDR", ":8080")
        fmt.Printf("serving %s at http://localhost%s\n", outputDir, addr)
        http.Handle("/", http.FileServer(http.Dir(outputDir)))
        if err := http.ListenAndServe(addr, nil); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }

    default:
        fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}

func getenv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
