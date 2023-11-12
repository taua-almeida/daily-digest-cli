package internal

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintPullRequests(detailedPullRequests []DetailedPullRequest, printConfig PrintStyle) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"PR Number", "Title", "PR Link", "Status", "PR Condition", "Mergeable", "CI/CD Status"})
	for _, prItem := range detailedPullRequests {
		t.AppendRow(table.Row{prItem.Number, prItem.Title, prItem.URL, prItem.State, prItem.Condition, prItem.IsMergeable, prItem.CICDStatus})
	}

	t.AppendFooter(table.Row{"Total", len(detailedPullRequests), "", "", ""})
	t.SetStyle(getStyle(printConfig.Style))
	t.Render()
}

func getStyle(styleName string) table.Style {
	switch styleName {
	case "bold":
		return table.StyleBold
	case "colored_bright":
		return table.StyleColoredBright
	case "colored_dark":
		return table.StyleColoredDark
	case "double":
		return table.StyleDouble
	case "light":
		return table.StyleLight
	case "rounder":
		return table.StyleRounded
	default:
		return table.StyleDefault
	}
}
