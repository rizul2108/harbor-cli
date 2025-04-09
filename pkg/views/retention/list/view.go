// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"golang.org/x/term"
)

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func formatScopeSelectors(selectors map[string][]models.RetentionSelector) string {
	var result []string
	for k, v := range selectors {
		var values []string
		for _, s := range v {
			values = append(values, fmt.Sprintf("%v", s))
		}
		result = append(result, fmt.Sprintf("%s: [%s]", k, strings.Join(values, ", ")))
	}
	return strings.Join(result, "; ")
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 160
	}
	return width
}

func getAdjustedColumns() []table.Column {
	totalWidth := getTerminalWidth()

	columnWidths := []int{
		totalWidth / 15, // Action
		totalWidth / 15, // Disabled
		totalWidth / 8,  // Params
		totalWidth / 15, // Priority
		totalWidth / 6,  // Scope Selectors
		totalWidth / 6,  // Tag Selectors
		totalWidth / 10, // Template
	}

	return []table.Column{
		{Title: "Action", Width: columnWidths[0]},
		{Title: "Disabled", Width: columnWidths[1]},
		{Title: "Params", Width: columnWidths[2]},
		{Title: "Priority", Width: columnWidths[3]},
		{Title: "Scope Selectors", Width: columnWidths[4]},
		{Title: "Tag Selectors", Width: columnWidths[5]},
		{Title: "Template", Width: columnWidths[6]},
	}
}

func ListRetentionRules(rules []*models.RetentionRule) {
	var rows []table.Row
	columns := getAdjustedColumns()

	for _, rule := range rules {
		params := ""
		for k, v := range rule.Params {
			params += fmt.Sprintf("%s: %v, ", k, v)
		}

		scopeSelectors := formatScopeSelectors(rule.ScopeSelectors)

		var tagSelectors []string
		for _, ts := range rule.TagSelectors {
			tagSelectors = append(tagSelectors, fmt.Sprintf("%v", ts))
		}

		rows = append(rows, table.Row{
			truncateString(rule.Action, columns[0].Width),
			truncateString(fmt.Sprintf("%v", rule.Disabled), columns[1].Width),
			truncateString(params, columns[2].Width),
			truncateString(fmt.Sprintf("%d", rule.Priority), columns[3].Width),
			truncateString(scopeSelectors, columns[4].Width),
			truncateString(strings.Join(tagSelectors, ", "), columns[5].Width),
			truncateString(rule.Template, columns[6].Width),
		})
	}

	m := tablelist.NewModel(columns, rows, len(rows))

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
