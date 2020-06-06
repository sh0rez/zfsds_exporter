package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func GetProps(dataset string, props []string) (map[string]string, error) {
	cmd := exec.Command("zfs", "list",
		"-p", // exact numbers
		"-H", // no table header
		"-o", strings.Join(props, ","),
		dataset,
	)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	m := make(map[string]string)
	out := strings.TrimSpace(buf.String())
	for i, s := range strings.Fields(out) {
		m[props[i]] = s
	}

	return m, nil
}
