package main

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

func shellIntegration(c *cli.Context) error {
	tmpl, err := template.New("bash").Parse(bashTemplate)
	if err != nil {
		return err
	}

	bin, err := os.Executable()
	if err != nil {
		return err
	}

	tmpl.Execute(c.App.Writer, bin)
	return nil
}

func execOpen(c *cli.Context) error {
	file := c.String("with-shell-integration")
	if file == "" {
		return nil
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	c.App.Metadata["shell-commands-file"] = f
	return nil
}

func execClose(c *cli.Context) error {
	if f, ok := c.App.Metadata["shell-commands-file"].(*os.File); ok {
		return f.Close()
	}

	return nil
}

// exec appends a shell command to the file specified by --shell-commands
func exec(c *cli.Context, v ...string) {
	f, ok := c.App.Metadata["shell-commands-file"].(*os.File)
	if !ok {
		return
	}

	for _, a := range v {
		a = "'" + strings.Replace(a, "'", `'\''`, -1) + "' "
		if _, err := io.WriteString(f, a); err != nil {
			panic(err)
		}
	}
	if _, err := io.WriteString(f, "\n"); err != nil {
		panic(err)
	}
}

var bashTemplate = `
grit() {
    local file="$(mktemp)"
    trap "rm -f '$file'" EXIT
    "{{.}}" --with-shell-integration="$file" "$@" && source "$file"
    return $?
}

# Setup autocompletion using the real binary.
_grit_bash_autocomplete() {
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$(GRIT_COMP_WORDS="${COMP_WORDS[@]}" "{{.}}" ${COMP_WORDS[@]:1:$COMP_CWORD} --generate-bash-completion)
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}

complete -F _grit_bash_autocomplete grit
`
