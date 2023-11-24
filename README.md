# carapace-spec-man

[Spec](https://github.com/rsteube/carapace-spec) generation for manpages.

```yaml
# carapace-spec-man touch
name: touch
description: change file timestamps
flags:
    --date: parse STRING and use it instead of current time
    --help: display this help and exit
    --no-create: do not create any files
    --no-dereference: affect each symbolic link instead of any referenced file (useful only on systems that can change the timestamps of a symlink)
    --reference: use this file's times instead of current time
    --time: 'change the specified time: WORD is access, atime, or use: equivalent to -a WORD is modify or mtime: equivalent to -m'
    --version: output version information and exit
    -a: change only the access time
    -c: do not create any files
    -d: parse STRING and use it instead of current time
    -f: (ignored)
    -h: affect each symbolic link instead of any referenced file (useful only on systems that can change the timestamps of a symlink)
    -m: change only the modification time
    -r: use this file's times instead of current time
    -t: use [[CC]YY]MMDDhhmm[.ss] instead of current time
completion:
    positionalany:
        - $files
```

> [!IMPORTANT]
> Manpages are highly inconsistent so the results will contain errors.
> 
> Issues you will encounter:
> - parsing failing completely
> - parsing being stuck (`git` has this issue)
> - missing flags or subcommands
> - invalid subcommands (`-` in manpage name is assumed as subcommand delimiter)
> - description not truncated well
> - all flags are assumed boolean
> - shorthand flags aren't grouped with their longhand 
>
> It is recommended to prepare them manually for [carapace-parse] instead.

[carapace-parse]:https://github.com/rsteube/carapace-bin/tree/master/cmd/carapace-parse
