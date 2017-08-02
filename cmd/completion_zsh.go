package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"bytes"
)

// Taken from the "kubectl" source code [1]
//   [1] https://raw.githubusercontent.com/kubernetes/kubernetes/master/pkg/kubectl/cmd/completion.go
// Apache-2.0 licensed

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates auto-complete definitions for zsh",
	Long: "Generates auto-complete definitions for zsh",
	Run: func(cmd *cobra.Command, args []string) {
		out := os.Stdout
		zsh_initialization := `
__spacectl_get_teams() {
	local teams="$(spacectl t ls | awk 'NR > 1 { printf "%s ",$2 }')"
	COMPREPLY=( $(__spacectl_compgen -W "${teams}" -- "$cur" ) )
}

__spacectl_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__spacectl_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__spacectl_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__spacectl_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__spacectl_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__spacectl_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}

__spacectl_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__spacectl_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__spacectl_filedir() {
	local RET OLD_IFS w qw

	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi

	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"

	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"

	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__spacectl_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__spacectl_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
    	printf %q "$1"
    fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__spacectl_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__spacectl_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__spacectl_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__spacectl_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__spacectl_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__spacectl_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__spacectl_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__spacectl_type/g" \
	<<'BASH_COMPLETION_EOF'
`
		out.Write([]byte(zsh_initialization))

		buf := new(bytes.Buffer)
		RootCmd.GenBashCompletion(buf)
		out.Write(buf.Bytes())

		zsh_tail := `
BASH_COMPLETION_EOF
}

__spacectl_bash_source <(__spacectl_convert_bash_to_zsh)
`
		out.Write([]byte(zsh_tail))
	},
}

func init() {
	completionCmd.AddCommand(completionZshCmd)
}
