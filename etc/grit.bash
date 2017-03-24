# Store the location of the real grit binary.
GRIT_BIN="${GRIT_BIN:-$(which grit)}"

if [ -f "$GRIT_BIN" ]; then

    # Redefine grit as a shell function (as opposited to a regular binary), so
    # we can affect the current shell, such as changing the working directory.
    grit() {
        # Create a tempory for file shell commands. Grit writes any commands
        # (such as directory changes) to this file.
        local file="$(mktemp)"
        trap "rm -f '$file'" EXIT

        # Run grit and execute the shell commands.
        $GRIT_BIN --shell-commands="$file" "$@" && source "$file"
        return $?
    }

    # Setup autocompletion using the real binary.
    _grit_bash_autocomplete() {
        local cur opts base
        COMPREPLY=()
        cur="${COMP_WORDS[COMP_CWORD]}"
        opts=$($GRIT_BIN ${COMP_WORDS[@]:1:$COMP_CWORD} --generate-bash-completion)
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    }

    complete -F _grit_bash_autocomplete grit

fi
