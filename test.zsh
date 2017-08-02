
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
# bash completion for spacectl                             -*- shell-script -*-

__debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__my_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__handle_reply()
{
    __debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%%=*}"
                __index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zfs completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        declare -F __custom_func >/dev/null && __custom_func
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__handle_flag()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    if [ -n "${flagvalue}" ] ; then
        flaghash[${flagname}]=${flagvalue}
    elif [ -n "${words[ $((c+1)) ]}" ] ; then
        flaghash[${flagname}]=${words[ $((c+1)) ]}
    else
        flaghash[${flagname}]="true" # pad "true" for bool flag
    fi

    # skip the argument to a two word flag
    if __contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__handle_noun()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__handle_command()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_$(basename "${words[c]//:/__}")"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__handle_word()
{
    if [[ $c -ge $cword ]]; then
        __handle_reply
        return
    fi
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __handle_flag
    elif __contains_word "${words[c]}" "${commands[@]}"; then
        __handle_command
    elif [[ $c -eq 0 ]] && __contains_word "$(basename "${words[c]}")" "${commands[@]}"; then
        __handle_command
    else
        __handle_noun
    fi
    __handle_word
}

_spacectl_apply()
{
    last_command="spacectl_apply"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_completion_bash()
{
    last_command="spacectl_completion_bash"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_completion_zsh()
{
    last_command="spacectl_completion_zsh"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_completion()
{
    last_command="spacectl_completion"
    commands=()
    commands+=("bash")
    commands+=("zsh")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_invites_accept()
{
    last_command="spacectl_invites_accept"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--invite=")
    two_word_flags+=("-i")
    local_nonpersistent_flags+=("--invite=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_flag+=("--invite=")
    must_have_one_flag+=("-i")
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_invites_list()
{
    last_command="spacectl_invites_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--out")
    local_nonpersistent_flags+=("--out")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_invites_revoke()
{
    last_command="spacectl_invites_revoke"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--invite=")
    two_word_flags+=("-i")
    local_nonpersistent_flags+=("--invite=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_invites()
{
    last_command="spacectl_invites"
    commands=()
    commands+=("accept")
    commands+=("list")
    commands+=("revoke")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_login()
{
    last_command="spacectl_login"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--auth-server=")
    local_nonpersistent_flags+=("--auth-server=")
    flags+=("--password=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--password=")
    flags+=("--username=")
    two_word_flags+=("-u")
    local_nonpersistent_flags+=("--username=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_apply()
{
    last_command="spacectl_apply"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces_delete()
{
    last_command="spacectl_spaces_delete"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--yes")
    flags+=("-y")
    local_nonpersistent_flags+=("--yes")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces_init()
{
    last_command="spacectl_spaces_init"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dns-label=")
    two_word_flags+=("-l")
    local_nonpersistent_flags+=("--dns-label=")
    flags+=("--force")
    local_nonpersistent_flags+=("--force")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_flag+=("--name=")
    must_have_one_flag+=("-n")
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces_list()
{
    last_command="spacectl_spaces_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces_open()
{
    last_command="spacectl_spaces_open"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--stage=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--stage=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces_show()
{
    last_command="spacectl_spaces_show"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces_validate()
{
    last_command="spacectl_spaces_validate"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_spaces()
{
    last_command="spacectl_spaces"
    commands=()
    commands+=("apply")
    commands+=("delete")
    commands+=("init")
    commands+=("list")
    commands+=("open")
    commands+=("show")
    commands+=("validate")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--spacefile=")
    two_word_flags+=("-f")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_teams_create()
{
    last_command="spacectl_teams_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dns-label=")
    two_word_flags+=("-l")
    local_nonpersistent_flags+=("--dns-label=")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_teams_invite()
{
    last_command="spacectl_teams_invite"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--email=")
    two_word_flags+=("-e")
    local_nonpersistent_flags+=("--email=")
    flags+=("--message=")
    two_word_flags+=("-m")
    local_nonpersistent_flags+=("--message=")
    flags+=("--role=")
    two_word_flags+=("-r")
    local_nonpersistent_flags+=("--role=")
    flags+=("--user-id=")
    two_word_flags+=("-u")
    local_nonpersistent_flags+=("--user-id=")
    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_teams_list()
{
    last_command="spacectl_teams_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_teams_members()
{
    last_command="spacectl_teams_members"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_teams()
{
    last_command="spacectl_teams"
    commands=()
    commands+=("create")
    commands+=("invite")
    commands+=("list")
    commands+=("members")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl_version()
{
    last_command="spacectl_version"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_spacectl()
{
    last_command="spacectl"
    commands=()
    commands+=("apply")
    commands+=("completion")
    commands+=("invites")
    commands+=("login")
    commands+=("spaces")
    commands+=("teams")
    commands+=("version")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--api-server=")
    flags+=("--config=")
    flags+=("--non-interactive")
    flags+=("--team=")
    flags_with_completion+=("--team")
    flags_completion+=("__spacectl_get_teams")
    two_word_flags+=("-t")
    flags_with_completion+=("-t")
    flags_completion+=("__spacectl_get_teams")
    flags+=("--token-file=")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_spacectl()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __my_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("spacectl")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_spacectl spacectl
else
    complete -o default -o nospace -F __start_spacectl spacectl
fi

# ex: ts=4 sw=4 et filetype=sh

BASH_COMPLETION_EOF
}

__spacectl_bash_source <(__spacectl_convert_bash_to_zsh)
