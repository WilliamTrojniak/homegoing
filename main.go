package main

import (
	"fmt"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
)

var myApp *App;

func main() {
  wd, err := os.Getwd();
  if err != nil {
    fmt.Printf("Error while getting current working directory: %v", err);
    os.Exit(1);
  }

  myApp = newApp(path.Join(wd, "dotfiles.toml"));
  p := tea.NewProgram(myApp, tea.WithAltScreen());
  if _, err := p.Run(); err != nil {
    fmt.Printf("Unexpected error occured: %v", err);
    os.Exit(1);
  }
}

