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
go get github.com/seanbreckenridge/spkglist/cmd/spkglist
```

```
>>>PMARK
perl -E 'print "`"x3, "\n"'
go build ./cmd/spkglist/
./spkglist -help
perl -E 'print "`"x3, "\n"'
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

For more examples, you can see my usage in my system bootstrap script [here](https://github.com/seanbreckenridge/dotfiles/blob/7c570944b244986d2837ffa935ff8efd7e7f4543/.config/yadm/computer_bootstrap#L37-L103), corresponding package lists [here](https://github.com/seanbreckenridge/dotfiles/tree/baf92d5fed00b87167b509f22d439c5e2075f63b/.config/yadm/package_lists).

## Specification

```
>>>PMARK
perl -E 'print "`"x3, "conf", "\n"'
cat ./examples/spec.txt
perl -E 'print "`"x3, "\n"'
```

Output:

```
>>>PMARK
perl -E 'print "`"x3, "\n"'
go build ./cmd/spkglist/
./spkglist -skip-last-delim ./examples/spec.txt 2>/dev/null
perl -E 'print "\n", "`"x3, "\n"'
```
