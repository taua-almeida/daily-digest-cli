package internal

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintPullRequests(detailedPullRequests []DetailedPullRequest) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"PR Number", "Title", "PR Link", "Status", "PR Condition", "Mergeable", "CI/CD Status"})
	for _, prItem := range detailedPullRequests {
		t.AppendRow(table.Row{prItem.Number, prItem.Title, prItem.URL, prItem.State, prItem.Condition, prItem.IsMergeable, prItem.CICDStatus})
	}

	t.AppendFooter(table.Row{"Total", len(detailedPullRequests), "", "", ""})
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}
