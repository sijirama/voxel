package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/jroimartin/gocui"
)

// fetchJournalctlLogs runs the journalctl command and returns the logs as a string
func fetchJournalctlLogs() (string, error) {
	cmd := exec.Command("journalctl", "-n", "50") // Fetch the last 50 logs
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// layout sets up the initial layout of the UI
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		logs, err := fetchJournalctlLogs()
		if err != nil {
			return fmt.Errorf("could not fetch logs: %v", err)
		}
		v.Write([]byte(logs))
	}
	if v, err := g.SetView("status", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "Press Ctrl+C to quit")
	}
	return nil
}

// keybindings sets up the keybindings for the application
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

// quit exits the application
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
