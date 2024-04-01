package cliutils

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	errorHeader   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#FF5F87")).Bold(true).Padding(0, 1).Margin(1).MarginLeft(2).SetString(" ERROR ")
	warningHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#FFAF5F")).Bold(true).Padding(0, 1).Margin(1).MarginLeft(2).SetString(" WARNING ")
	successHeader = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#00BFA5")).Bold(true).Padding(0, 1).Margin(1).MarginLeft(2).SetString(" SUCCESS ")

	messageDetails = lipgloss.NewStyle().Foreground(lipgloss.Color("#757575")).Margin(0, 0, 1, 2)
)

func PrintError(title string, err error) {
	fmt.Printf("%s\n", lipgloss.JoinHorizontal(lipgloss.Center, errorHeader.String(), title))
	fmt.Printf("%s\n", messageDetails.Render(err.Error()))
}

func PrintErrorFatal(title string, err error) {
	PrintError(title, err)
	os.Exit(1)
}

func PrintWarning(title string, message string) {
	fmt.Printf("%s\n", lipgloss.JoinHorizontal(lipgloss.Center, warningHeader.String(), title))
	fmt.Printf("%s\n", messageDetails.Render(message))
}

func PrintSuccess(title string, message string) {
	fmt.Printf("%s\n", lipgloss.JoinHorizontal(lipgloss.Center, successHeader.String(), title))
	fmt.Printf("%s\n", messageDetails.Render(message))
}
