# Store the location of the real grit binary.
GRIT_BIN=$(which grit)

if [ -f "$GRIT_BIN" ]; then

    # Make a bash function named grit so that we can execute our directory changes
    # in the current shell.
    grit() {
        case $1 in
            # Intercept commands that print directories to STDOUT.
            clone|cd)
                local dir=$($GRIT_BIN "$@")
                if [ -d "$dir" ]; then
                    cd "$dir"
                    return 0
                fi

                return 1
            ;;
            # Pass through other commands unchanged.
            *)
                $GRIT_BIN "$@"
                return $?
            ;;
        esac
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
