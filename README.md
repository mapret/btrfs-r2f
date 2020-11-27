# btrfs-r2f

An application to receive a btrfs command stream to a local folder on any file system.

## Goal

> Btrfs is a modern copy on write (CoW) filesystem for Linux aimed at implementing advanced features while also focusing on fault tolerance, repair and easy administration.  
>&nbsp;&nbsp;&nbsp;&nbsp;— [btrfs.wiki.kernel.org](https://btrfs.wiki.kernel.org)

Btrfs can make [atomic snapshots](https://man7.org/linux/man-pages/man8/btrfs-subvolume.8.html) and store them as [files](https://man7.org/linux/man-pages/man8/btrfs-send.8.html).
The additionally functionality of storing only the differences to a previous snapshot makes this mechanism suitable for [incremental backups](https://en.wikipedia.org/wiki/Incremental_backup).
Restoring a backup from these files requires a btrfs partition, where all increments [are unpacked](https://man7.org/linux/man-pages/man8/btrfs-receive.8.html) in the same order as they were created.

Requiring a btrfs partition to restore a backup is no major problem if the backup is used for restoring the filesystem on a server, which would need the partition anyways.
Testing a backup, on the other hand, might be a slight inconvenience, since most Linux distributions do not use btrfs [by default](https://btrfs.wiki.kernel.org/index.php/Production_Users). 
Therefore, this tool was created.
It can unpack (=receive) a btrfs command stream to a directory on any filesystem, which makes testing a sequence of incremental backups more accessible.
Hence the name, r2f, receive to folder.

## Building

This project is implemented in the [Go programming language](https://golang.org/).

To build the application, the [Go compiler](https://golang.org/doc/install) needs to be installed and available to the commandline.
From the root of the repository, run the command
```
go build
```
, which generates the executable.

## Running

```
Usage: btrfs-r2f [options...]

-i --input
    Command stream file, stdin is used if not supplied
-o --output
    Target directory, default value "btrfs-subvolume"
-d --dry-run
    Dont actually create files, only print commands (useful with -v)
-v --verbose
    Print every command that is processed to stdout
```

## Testing

A suite of simple tests can be run with the built-in command ```go test```.

Additionally, exploratory tests are run from time to time, where a series of incremental backups is unpacked and then compared to the original snapshot.
This is done by calculating the hash of each file in a btrfs snapthos with the command ```find . -type f -print0 | sort -z | xargs -r0 sha256sum```,
then running the same command on the restored directory and comparing the outputs.
As of writing this, the exploratory test succeeded for a partition containing about 21000 files with a total size of about 20GiB, contained in 19 backup files (full initial backup plus daily incremental backup). 

## Remarks

On Windows, symlinks cannot be created without administrator privileges by default.
Therefore, symlinks are represented as shortcuts.
This has the disadvantage that restored backups are not exactly equal to their original source, because shortcuts always end in .lnk, which is not shown in Windows Explorer, but becomes relevant when reading file names through a program.

## License


This project is licensed under the **GPLv3**.
See [LICENSE.txt](LICENSE.txt) for the full license text.
Since the short license notice recommended by the GPLv3 is not present in the source files, it is displayed here instead:

```
btrfs-r2f, an application to receive a btrfs command stream to a
local folder on any file system.
Copyright (C) 2020  Markus Prettner

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
