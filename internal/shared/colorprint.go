package shared

import "github.com/fatih/color"

var (
	Red     = color.New(color.FgRed, color.Bold).PrintfFunc()
	Sred    = color.New(color.FgRed, color.Bold).SprintfFunc()
	Magenta = color.New(color.FgMagenta, color.Bold).PrintfFunc()
	Green   = color.New(color.FgGreen, color.Bold).PrintfFunc()
	Faint   = color.New(color.Faint).PrintfFunc()
	Sfaint  = color.New(color.Faint).SprintfFunc()
	Cyan    = color.New(color.FgCyan).PrintfFunc()
	Yellow  = color.New(color.FgYellow).PrintfFunc()
	Syellow = color.New(color.FgYellow).SprintfFunc()

	Colors = []string{"gold", "firebrick", "green", "darksalmon", "aquamarine", "moccasin", "turquoise"}
)
