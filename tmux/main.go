package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var count = 0
var session_name = "DevMenu"

// Runs a tmux command with args in the shell
func RunTmuxCmd(args []string) (string, string, error) {
	cmd := exec.Command("tmux", args...)

	fmt.Printf("%s\n", cmd.String())
	log := exec.Command("echo", string(count))
	err := log.Run()

	count++

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	return outStr, errStr, err
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func (w Direction) String() string {
	return [...]string{"Horizontal", "Vertical"}[w]
}

func (w Direction) EnumIndex() int {
	return int(w)
}

type Row struct {
	paneId string
}

type Column struct {
	paneId   string
	children []*Row
}

var columns []*Column = []*Column{{
	paneId: "",
	children: []*Row{{
		paneId: "",
	}, {
		paneId: "",
	}, {
		paneId: "",
	}},
}, {
	paneId: "",
	children: []*Row{{
		paneId: "",
	}, {
		paneId: "",
	}},
}, {
	paneId: "",
	children: []*Row{{
		paneId: "",
	}, {
		paneId: "",
	}, {
		paneId: "",
	}, {
		paneId: "",
	}},
}}

func splitWindow(direction Direction, target string) (string, error) {
	fmt.Printf("Splitting %s target: %s\n", direction, target)

	var cmdArgs []string
	if direction == Horizontal {
		cmdArgs = []string{"split-window", "-h", "-P", "-F", "#{pane_id}", "-t", target}
	} else {
		cmdArgs = []string{"split-window", "-P", "-F", "#{pane_id}", "-t", target}
	}

	out, errStr, err := RunTmuxCmd(cmdArgs)
	if err != nil {
		return "", fmt.Errorf("split-window failed: %v (%s)", err, errStr)
	}

	// tmux returns the new pane_id in `out`.
	newPaneId := strings.TrimSpace(out)
	return newPaneId, nil

}

func getInitialPaneId() (string, error) {
	cmd := exec.Command("bash", "-c", "tmux list-panes -t DevMenu -F '#{pane_id}' | head -n 1")
	outBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	paneId := strings.TrimSpace(string(outBytes))

	// TODO: Verbose
	fmt.Printf("Found first pane %s\n", paneId)

	return paneId, nil
}

func renderColumns() {
	for i, column := range columns {
		// If it's the first pane, set the paneId to the initialPaneId
		// returned from tmux
		if i == 0 {
			initialPaneId, err := getInitialPaneId()

			if err != nil {
				panic(err)
			}

			column.paneId = initialPaneId
		}

		hasNextColumn := i < (len(columns) - 1)

		// We only want to split the current column if there is one next
		if hasNextColumn {
			newPaneId, err := splitWindow(Horizontal, column.paneId)

			if err != nil {
				panic(err)
			}

			// Set the next columns paneId
			columns[i+1].paneId = newPaneId
		}

	}
}

func resizeRowsInColumn(column *Column) {
	if len(column.children) == 0 {
		return
	}

	height, err := getWindowHeight()

	if err != nil {
		panic(err)
	}

	rows := len(column.children)

	rowHeight := height / rows

	for _, row := range column.children {
		_, _, err := RunTmuxCmd([]string{"resize-pane", "-t", row.paneId, "-y", strconv.Itoa(rowHeight)})

		if err != nil {
			panic(err)
		}
	}
}

func renderRows() {
	for _, column := range columns {
		if len(column.children) > 0 {
			for r, row := range column.children {
				target := row.paneId
				hasNextRow := r < (len(column.children) - 1)

				// If this is the first row, target the parent column
				if r == 0 {
					row.paneId = column.paneId
					target = column.paneId
				}

				if hasNextRow {
					// Only split vertically when there is a next row
					nextPaneId, err := splitWindow(Vertical, target)

					if err != nil {
						panic(err)
					}

					// Set the next rows paneId
					column.children[r+1].paneId = nextPaneId
				}

			}

		}

		resizeRowsInColumn(column)
	}

}

func getWindowHeight() (int, error) {
	heightStr, _, err := RunTmuxCmd([]string{"display-message", "-p", "#{window_height}"})

	if err != nil {
		return 0, err
	}

	clean := strings.TrimSpace(heightStr)

	height, err := strconv.Atoi(clean)

	if err != nil {
		return 0, err
	}

	return height, nil
}

func main() {
	renderColumns()
	// Now the first set of columns is rendered, select the layout
	RunTmuxCmd([]string{"select-layout", "-n"})
	renderRows()
}
