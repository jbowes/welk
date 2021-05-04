package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jbowes/sumdog/internal/install"

	"github.com/charmbracelet/bubbles/progress"
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

		prog, err := progress.NewModel(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
		if err != nil {
			fmt.Println("Could not initialize progress model:", err)
			os.Exit(1)
		}

		if err = tea.NewProgram(example{progress: prog}).Start(); err != nil {
			fmt.Println("Oh no!", err)
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

const (
	fps              = 60
	stepSize float64 = 1.0 / (float64(fps) * 2.0)
	padding          = 2
	maxWidth         = 80
)

type tickMsg time.Time

type example struct {
	percent  float64
	progress *progress.Model
}

func (e example) Init() tea.Cmd {
	return tickCmd()
}

func (e example) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return e, tea.Quit
		default:
			return e, nil
		}

	case tea.WindowSizeMsg:
		e.progress.Width = msg.Width - padding*2 - 4
		if e.progress.Width > maxWidth {
			e.progress.Width = maxWidth
		}
		return e, nil

	case tickMsg:
		e.percent += stepSize
		if e.percent > 1.0 {
			e.percent = 1.0
			return e, tea.Quit
		}
		return e, tickCmd()

	default:
		return e, nil
	}
}

func (e example) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" + pad + e.progress.View(e.percent) + "\n\n"
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
