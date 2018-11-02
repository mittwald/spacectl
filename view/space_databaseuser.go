package view

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/cmd/helper"
	"io"
	"time"
)

type DatabaseUserView interface {
	List(userList spaces.DatabaseUserList, out io.Writer)
}

type TabularDatabaseUserView struct{}

func (t TabularDatabaseUserView) List(userList spaces.DatabaseUserList, out io.Writer) {
	if len(userList) == 0 {
		fmt.Fprintln(out, "No users found.")
		return
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("USER", "STATUS", "TYPE", "CREATED")

	for _, user := range userList {
		since := helper.HumanReadableDateDiff(time.Now(), user.CreatedAt)

		table.AddRow(
			user.User,
			user.Status,
			user.Type,
			since + " ago",
		)
	}

	fmt.Fprintln(out, table)
}
