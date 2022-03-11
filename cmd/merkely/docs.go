package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const docsDesc = `
Generate documentation files for Merkely CLI.
This command can generate documentation in the following formats: Markdown.
`

type docsOptions struct {
	dest            string
	topCmd          *cobra.Command
	generateHeaders bool
}

func newDocsCmd(out io.Writer) *cobra.Command {
	o := &docsOptions{}

	cmd := &cobra.Command{
		Use:    "docs",
		Short:  "generate documentation as markdown",
		Long:   docsDesc,
		Hidden: true,
		Args:   NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.topCmd = cmd.Root()
			return o.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&o.dest, "dir", "./", "The directory to which documentation is written.")
	f.BoolVar(&o.generateHeaders, "generate-headers", true, "Generate standard headers for markdown files.")

	return cmd
}

func (o *docsOptions) run() error {
	if o.generateHeaders {
		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/client_reference/" + strings.ToLower(base) + "/"
		}

		hdrFunc := func(filename string) string {
			base := filepath.Base(filename)
			name := strings.TrimSuffix(base, path.Ext(base))
			title := strings.ToLower(strings.Replace(name, "_", " ", -1))
			return fmt.Sprintf("---\ntitle: \"%s\"\n---\n\n", title)
		}

		return MereklyGenMarkdownTreeCustom(o.topCmd, o.dest, hdrFunc, linkHandler)
	}
	return doc.GenMarkdownTree(o.topCmd, o.dest)
}

func MereklyGenMarkdownTreeCustom(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := MereklyGenMarkdownTreeCustom(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}

	if !cmd.HasParent() || !cmd.HasSubCommands() {
		basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"
		filename := filepath.Join(dir, basename)
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
			return err
		}
		if err := MerkelyGenMarkdownCustom(cmd, f, linkHandler); err != nil {
			return err
		}
	}
	return nil
}

// MerkelyGenMarkdownCustom creates custom markdown output.
func MerkelyGenMarkdownCustom(cmd *cobra.Command, w io.Writer, linkHandler func(string) string) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("## " + name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```shell\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```shell\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd, name); err != nil {
		return err
	}
	// if !cmd.DisableAutoGenTag {
	// 	buf.WriteString("###### Auto generated by spf13/cobra on " + time.Now().Format("2-Jan-2006") + "\n")
	// }
	_, err := buf.WriteTo(w)
	return err
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command, name string) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("### Options\n")
		buf.WriteString("| Flag | Description |\n")
		buf.WriteString("| :--- | :--- |\n")
		usages := CommandsInTable(flags)
		fmt.Fprint(buf, usages)
		buf.WriteString("\n\n")
	}

	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Options inherited from parent commands\n")
		buf.WriteString("| Flag | Description |\n")
		buf.WriteString("| :--- | :--- |\n")
		usages := CommandsInTable(parentFlags)
		fmt.Fprint(buf, usages)
		buf.WriteString("\n\n")
	}
	return nil
}
