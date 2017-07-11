Peco Actions
============

This command provides the features:

* Show an appropriate menu for input type.
* Choose an action from the menu.


## Environments

* Only macOS so far.


## Installation

```
$ go get github.com/mattn/go-shellwords
$ make
$ cp peco-actions /usr/local/bin
```


## Usage

### usage

```sh
usage: peco-actions [option]
options:
  -h,--help     show this usage
  -v,--version  display the version
  -D,--debug    display the version
  --file        actions for file path(s)
  --process     actions for a process id
  --server      actions for a host/IP-address
  --git         actions for a commit id
  --generic     actions for generic only using addhoc menu
```


### zsh

```sh
__peco_actions__run_action() {
    local action="$*"
    if [ -n "$action" ]; then
        echo "$fg_bold[yellow]>> $fg_bold[cyan]$action$reset_color" >&2

        # Add command history
        print -s $action

        # Run by eval to handle arguments including spaces correctly
        # http://labs.opentone.co.jp/?p=5651
        eval $action
    fi
}
```

```sh
$ __peco_actions__run_action $(ls -1 | peco-actions --file)
$ __peco_actions__run_action $(ps | peco | awk '{ print $1 }' | peco-actions --process)
```
