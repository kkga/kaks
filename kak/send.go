package kak

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Send(kctx *Context, kakCommand string, errOutFile *os.File) error {
	kakExec, err := kakExec()
	if err != nil {
		return err
	}
	cmd := exec.Command(kakExec, "-p", kctx.Session.Name)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	var sb strings.Builder

	// wrap Kakoune command in try-catch
	// try
	sb.WriteString("try %{")
	sb.WriteString(" eval")
	if kctx.Buffer.Name != "" {
		sb.WriteString(fmt.Sprintf(" -buffer %s", strings.ReplaceAll(kctx.Buffer.Name, " ", "\\ ")))
	} else if kctx.Client.Name != "" {
		sb.WriteString(fmt.Sprintf(" -try-client %s", kctx.Client.Name))
	}
	sb.WriteString(fmt.Sprintf(" %s", kakCommand))
	sb.WriteString(" }")

	// catch
	sb.WriteString(" catch %{")
	// echo error to Kakoune's debug buffer
	sb.WriteString(" echo -debug kks: %val{error}\n")
	if errOutFile != nil {
		// write a prefixed error to tmp file so that we can parse it in runner and decide what to do
		sb.WriteString(fmt.Sprintf(" echo -to-file %s %s %%val{error}", errOutFile.Name(), EchoErrPrefix))
		sb.WriteString("\n")
	}
	// echo error in client
	sb.WriteString(" eval")
	if kctx.Client.Name != "" {
		sb.WriteString(fmt.Sprintf(" -try-client %s", kctx.Client.Name))
	}
	sb.WriteString(" %{ echo -markup {Error}kks: %val{error} }")
	sb.WriteString(" }")

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, sb.String()) //nolint
	}()

	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
