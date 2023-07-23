package util

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

func ConcatVertically(lhs, rhs string) string {
	builder := strings.Builder{}
	table := tabwriter.NewWriter(&builder, 1, 4, 1, ' ', 0)
	lhs_split := strings.Split(lhs, "\n")
	rhs_split := strings.Split(rhs, "\n")
	maxlen := Max(len(lhs_split), len(rhs_split))

	for i := 0; i < maxlen; i++ {
		if i >= len(lhs_split) {
			fmt.Fprint(table, "...")
		} else {
			fmt.Fprint(table, lhs_split[i])
		}
		fmt.Fprint(table, "\t#\t")
		if i >= len(rhs_split) {
			fmt.Fprint(table, "...")
		} else {
			fmt.Fprint(table, rhs_split[i])
		}
		fmt.Fprintln(table, "")
	}
	table.Flush()
	return builder.String()
}
