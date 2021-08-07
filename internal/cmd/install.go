// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jbowes/welk/internal/install"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/spf13/cobra"
)

var mu sync.Mutex // TODO: no global here.

var installCmd = cobra.Command{
	Use:     "install URL",
	Short:   "Safely install curl | sh style scripts",
	Aliases: []string{"i"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running install")

		err := install.Run(context.TODO(), permittedExec, log, args[0])
		if err != nil {
			fmt.Println("error installing", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(&installCmd)
}

func permittedExec(args []string) bool {
	mu.Lock()
	defer mu.Unlock()

	c := confirm{
		message: "Run external command: " + strings.Join(args, " "),
	}

	if err := tea.NewProgram(&c).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

	return *c.choice
}

func log(tag string, msg ...string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9999")).Render(tag) + " " + strings.Join(msg, " "))
}

type confirm struct {
	message string
	choice  *bool
}

func (c *confirm) Init() tea.Cmd { return nil }

func Noop() tea.Msg { return nil }

func (c *confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c": // TODO: handle this properly
			return c, tea.Quit
		case "enter", "n", "N":
			f := false
			c.choice = &f
			return c, tea.Sequentially(Noop, tea.Quit)
		case "y", "Y":
			t := true
			c.choice = &t
			return c, tea.Sequentially(Noop, tea.Quit)
		}
	}

	return c, nil
}

// TODO: needs a cursor
func (c *confirm) View() string {
	choice := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("n")
	if c.choice != nil {
		if *c.choice {
			choice = "y\n"
		} else {
			choice = "n\n"
		}

	}
	return c.message + "? " + choice
}
