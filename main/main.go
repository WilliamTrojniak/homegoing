package main

import (
	"fmt"
	"os"

  tea "github.com/charmbracelet/bubbletea"
)

var myApp *app;

func main() {
  myApp = newApp();
  p := tea.NewProgram(myApp);
  if _, err := p.Run(); err != nil {
    fmt.Printf("Unexpected error occured: %v", err);
    os.Exit(1);
  }
}

