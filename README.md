## spkglist

A simple package list parser.

This isn't for any language/tool in particular, its meant to clean up hand written package files for any tool, especially while creating scripts to bootstrap installation processes.

For example, if you had a list of system packages you wanted to install:

`packages.txt`:

```conf
# packages I want to install with my package manager
java  # ....
go
  python3
```

This removes comments and cleans up whitespace/newlines, so that that doesn't have to be done in the shell.

`$ spkglist packages.txt`

```conf
java
go
python3
```

### Install

```
go get github.com/purarue/spkglist/cmd/spkglist
```

```
A simple package list parser
Pass one or more files to read from as arguments
If none provided, reads from STDIN.
  -delimiter string
    	delimiter to print between results (default "\n")
  -json
    	print results as a JSON array
  -print0
    	use the null character as the delimiter
  -skip-last-delim
    	dont print the delimiter after the last item
  -split
    	split on all whitespace, instead of newlines
```

### Common Usage

To pass the parsed information to a command, one could do:

```
spkglist -print0 packages.txt | xargs -0 sudo apt install
```

Some `bash` that might be used to check if packages are installed, else install them:

```bash
# create a variable list of installed packages
CARGO_INSTALLED_PACKAGES="$(cargo install --list | sed -E -e '/^\s+/d; s|\s.*||')"
while read -r cargopkg; do
	# if we can't find that package in the installed packages
	if ! grep -q "^${cargopkg}$" <<<"${CARGO_INSTALLED_PACKAGES}"; then
		printf "Installing %s\n" "${cargopkg}"
		cargo install "${cargopkg}"
	fi
done < <(spkglist /path/to/package/list)
```

Or, you can query the package manager itself

```bash
# have to use for loop, while loop times out instantly
# when trying to prompt
#
# for complications with prompting while looping, see
# https://stackoverflow.com/q/6883363/9348376
# https://github.com/koalaman/shellcheck/wiki/SC2013
# http://mywiki.wooledge.org/DontReadLinesWithFor
for lib in $(spkglist /path/to/package/list); do
	if [[ ! $(yay -Q "${lib}" 2>/dev/null) ]]; then # if package isn't installed
		yay -S "${lib}"
	fi
done
```

For more examples, you can see my usage in my system bootstrap script [here](https://github.com/purarue/dotfiles/blob/7c570944b244986d2837ffa935ff8efd7e7f4543/.config/yadm/computer_bootstrap#L37-L103), corresponding package lists [here](https://github.com/purarue/dotfiles/tree/baf92d5fed00b87167b509f22d439c5e2075f63b/.config/yadm/package_lists).

## Specification

```conf
go  # just a package name

package name # assumes you want 'package name', you can pass the -split flag otherwise

# this would be an error, not allowed as a 'bare' character
😀

# you can quote any line with backticks, to include any character, including a '#'
`github.com/user/emoji_package😀#master`  # install from github

# python/node-esque packages work fine as bare words
requests==2.24.0
@elm-tooling/elm-language-server
serve
```

Output:

```
go
package name
github.com/user/emoji_package😀#master
requests==2.24.0
@elm-tooling/elm-language-server
serve
```

