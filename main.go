package main

import (
    "fmt"
    "net/http"
    "os"
    "path"
    "path/filepath"
    "strings"

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
        http.Handle("/", staticHandler(outputDir))
        if err := http.ListenAndServe(addr, nil); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }

    default:
        fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}

func staticHandler(outputDir string) http.Handler {
    fileServer := http.FileServer(http.Dir(outputDir))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if servesExistingPath(outputDir, r.URL.Path) {
            fileServer.ServeHTTP(w, r)
            return
        }

        notFoundPath := filepath.Join(outputDir, "404.html")
        html, err := os.ReadFile(notFoundPath)
        if err != nil {
            http.NotFound(w, r)
            return
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusNotFound)
        _, _ = w.Write(html)
    })
}

func servesExistingPath(outputDir, urlPath string) bool {
    cleanPath := path.Clean("/" + urlPath)
    relPath := strings.TrimPrefix(cleanPath, "/")
    filePath := filepath.Join(outputDir, relPath)

    info, err := os.Stat(filePath)
    if err != nil {
        return false
    }

    if !info.IsDir() {
        return true
    }

    _, err = os.Stat(filepath.Join(filePath, "index.html"))
    return err == nil
}

func getenv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
