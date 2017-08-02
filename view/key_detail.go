package view

import (
	"github.com/mittwald/spacectl/client/sshkeys"
	"io"
	"fmt"
	"github.com/gosuri/uitable"
	"time"
	"encoding/base64"
	"github.com/mittwald/spacectl/cmd/helper"
	"crypto/md5"
)

type KeyDetailView interface {
	KeyDetail(key *sshkeys.SSHKey, out io.Writer)
}

type TabularKeyDetailView struct {}

func (t TabularKeyDetailView) KeyDetail(key *sshkeys.SSHKey, out io.Writer) {
	fmt.Fprintln(out, "GENERAL INFO")

	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	since := helper.HumanReadableDateDiff(time.Now(), key.CreatedAt)

	keyStr := base64.StdEncoding.EncodeToString(key.Key)

	table.AddRow("  ID:", key.ID)
	table.AddRow("  Created:", since + " ago")
	table.AddRow("  Created At:", key.CreatedAt.String())
	table.AddRow("  Key:")
	table.AddRow("    Algorithm:", key.CipherAlgorithm)
	table.AddRow("    Comment:", key.Comment)
	table.AddRow("    Fingerprint:", fmt.Sprintf("%x", md5.Sum(key.Key)))

	fmt.Fprintln(out, table)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "PUBLIC KEY (BASE64)")
	fmt.Fprintln(out, keyStr)

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "PUBLIC KEY (BYTES)")

	i := 0
	for _, b := range key.Key {
		if (i % 32) == 0 {
			fmt.Fprint(out, "  ")
		}

		fmt.Fprintf(out, "%02x", b)

		if (i % 32) == 31 {
			fmt.Fprint(out, "\n")
		} else {
			fmt.Fprint(out, " ")
		}

		i ++
	}

	fmt.Fprintln(out, "")
}
