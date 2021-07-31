module github.com/jbowes/sumdog

go 1.15

require (
	github.com/benhoyt/goawk v1.7.0
	github.com/charmbracelet/bubbles v0.8.0
	github.com/charmbracelet/bubbletea v0.14.1
	github.com/charmbracelet/lipgloss v0.3.0
	github.com/containerd/console v1.0.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.9.0 // indirect
	github.com/rwtodd/Go.Sed v0.0.0-20190103233418-906bc69c9394
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.4.2
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	mvdan.cc/sh/v3 v3.2.4
)

// v3.3.0 doesn't work properly with the wasme install. remove exclude when
// it's past this version
exclude mvdan.cc/sh/v3 v3.3.0
