package main

import (
    "fmt"
    "os"
    "flag"
)

type Logger struct {
    Verbose bool
}

const LOG_FORMAT = "  > %s\n"

func NewLogger() *Logger {
    return &Logger{
        false,
    }
}

func (l *Logger) SetVerbose(v bool) {
    l.Verbose = v
}

func (l Logger) Fatal(line string, args... interface{}) {
    fLine := fmt.Sprintf(line, args...)
    fmt.Fprintf(os.Stderr, LOG_FORMAT, fLine)
    os.Exit(1)
}

func (l Logger) Info(line string, args... interface{}) {
    fLine := fmt.Sprintf(line, args...)
    fmt.Fprintf(os.Stdout, LOG_FORMAT, fLine)
    return
}

func (l Logger) Debug(line string, args... interface{}) {
    if l.Verbose {
        fLine := fmt.Sprintf(line, args...)
        fmt.Fprintf(os.Stdout, LOG_FORMAT, fLine)
    }
    return
}

func (l Logger) Usage() {
    fmt.Println("  Backup multiple MySQL hosts and Databases from one place\n")
    fmt.Println("Usage:")
    flag.PrintDefaults()
    os.Exit(0)
}