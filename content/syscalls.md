Title: An introduction to Linux system calls - and strace
Date: 2019-05-05
Category: Homepage
Tags: sysadmin, linux, system call, syscall, strace
Slug: syscall
Authors: hvnsweeting
Summary: see what's going on under the hood

## What is syscall - system call

Excerpt from `man 2 syscalls`

```
NAME
       syscalls - Linux system calls

SYNOPSIS
       Linux system calls.

DESCRIPTION
       The system call is the fundamental interface between an application and the Linux kernel.
```

Any (useful) program on GNU/Linux OS would need to ask Linux kernel do
something, e.g open/read/write to a file, use network, or memory...

```
PROGRAM <--------------> Linux kernel <---> hardware.
```

## How many syscalls?
There are ~ 403 syscall(s) as of Linux version

```
$ uname -r
4.15.0-46-generic
```

```
# get from man 2 syscalls output all lines which contain "(2)   number.number"
$ man 2 syscalls | grep -E '\(2\) +[0-9]\.[0-9]*' | wc -l
403
```

## Tools for observing syscall

`strace`

```
$ whatis strace
strace (1)           - trace system calls and signals
```

### Examples

#### How `free` command work?

Form: `strace command`

Output will be (very) verbose.

Let see if `free` cmd uses any `open*` syscall to open (to later read data
from) files:

```text
$ strace free 2>&1 | grep open
openat(AT_FDCWD, "/etc/ld.so.cache", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/usr/lib/x86_64-linux-gnu/libgtk3-nocsd.so.0", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libprocps.so.6", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libc.so.6", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libdl.so.2", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libpthread.so.0", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libsystemd.so.0", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/librt.so.1", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/liblzma.so.5", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/usr/lib/x86_64-linux-gnu/liblz4.so.1", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libgcrypt.so.20", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libgpg-error.so.0", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/proc/sys/kernel/osrelease", O_RDONLY) = 3
openat(AT_FDCWD, "/sys/devices/system/cpu/online", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/usr/lib/locale/locale-archive", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/proc/sys/kernel/osrelease", O_RDONLY) = 3
openat(AT_FDCWD, "/proc/meminfo", O_RDONLY) = 3
openat(AT_FDCWD, "/usr/share/locale/locale.alias", O_RDONLY|O_CLOEXEC) = 4
openat(AT_FDCWD, "/usr/share/locale/en/LC_MESSAGES/procps-ng.mo", O_RDONLY) = -1 ENOENT (No such file or directory)
openat(AT_FDCWD, "/usr/share/locale-langpack/en/LC_MESSAGES/procps-ng.mo", O_RDONLY) = -1 ENOENT (No such file or directory)
```

Skip all `.so` files - which stand for shared object in a dynamic library
that the program uses, and files that do not exist (ENOENT - see `man 3 errno`)

```text
$ strace free 2>&1 | grep open | grep -vF .so | grep -v ENOENT
openat(AT_FDCWD, "/proc/sys/kernel/osrelease", O_RDONLY) = 3
openat(AT_FDCWD, "/sys/devices/system/cpu/online", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/usr/lib/locale/locale-archive", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/proc/sys/kernel/osrelease", O_RDONLY) = 3
openat(AT_FDCWD, "/proc/meminfo", O_RDONLY) = 3
openat(AT_FDCWD, "/usr/share/locale/locale.alias", O_RDONLY|O_CLOEXEC) = 4
```

Turns out, `free` reads data from `/proc/meminfo` (which exposed by Linux kernel).

#### Where does `uptime` cmd get data from?

```
$ uptime
 00:39:59 up 2 days, 11:33,  3 users,  load average: 0.00, 0.04, 0.08
```

`strace` option `-y` print paths associated with file descriptor arguments.

```text
$ strace -y uptime 2>&1 | grep read | grep -vF .so
read(3</proc/sys/kernel/osrelease>, "4.15.0-46-generic\n", 1024) = 18
read(3</sys/devices/system/cpu/online>, "0-3\n", 8192) = 4
read(3</usr/share/zoneinfo/Asia/Ho_Chi_Minh>, "TZif2\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\6\0\0\0\6\0\0\0\0"..., 4096) = 389
read(3</usr/share/zoneinfo/Asia/Ho_Chi_Minh>, "TZif2\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\6\0\0\0\6\0\0\0\0"..., 4096) = 221
read(3</proc/uptime>, "214389.48 91487.77\n", 8191) = 19
read(4</run/utmp>, "\2\0\0\0\0\0\0\0~\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"..., 384) = 384
read(4</run/utmp>, "\6\0\0\0\252\4\0\0tty1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"..., 384) = 384
read(4</run/utmp>, "\7\0\0\0\374)\0\0tty7\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"..., 384) = 384
read(4</run/utmp>, "\1\0\0\0005\0\0\0~\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"..., 384) = 384
read(4</run/utmp>, "\7\0\0\0\321V\0\0pts/6\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"..., 384) = 384
read(4</run/utmp>, "\7\0\0\0\321V\0\0pts/9\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0"..., 384) = 384
read(4</run/utmp>, "", 384)             = 0
read(4</proc/loadavg>, "0.01 0.05 0.08 1/461 30551\n", 8191) = 27
```

### What are most common used syscalls?
It actually depends on which program, a command like `free` would not need
network access like `ping`.

`strace` option `-c` shows statistics, use `-S calls` to sorts by most called
syscall (desc).

```
$ strace -cS calls free 2>&1 | head
              total        used        free      shared  buff/cache   available
Mem:        3943388      842660      700960      109152     2399768     2704556
Swap:             0           0           0
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
  0.00    0.000000           0        33           mmap
  0.00    0.000000           0        24           mprotect
 53.45    0.000031           2        20         2 openat
 17.24    0.000010           1        19           close
  0.00    0.000000           0        17           read
```

For `ping`

```
$ sudo strace -cS calls ping -c1 1.1.1.1
PING 1.1.1.1 (1.1.1.1) 56(84) bytes of data.
64 bytes from 1.1.1.1: icmp_seq=1 ttl=58 time=32.4 ms

--- 1.1.1.1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 32.487/32.487/32.487/0.000 ms
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
 14.63    0.000067           4        16           mmap
 12.66    0.000058           5        12           mprotect
  2.18    0.000010           1         8           close
  3.28    0.000015           2         8           fstat
  4.80    0.000022           3         7         7 access
  3.49    0.000016           2         7           setsockopt
  1.97    0.000009           1         7           capget
  7.21    0.000033           5         7           openat
 17.47    0.000080          13         6           write
  2.62    0.000012           2         5           read
  4.59    0.000021           4         5         2 socket
  1.09    0.000005           2         3           brk
  0.66    0.000003           1         3           rt_sigaction
  1.31    0.000006           2         3           capset
  0.66    0.000003           2         2           ioctl
  0.44    0.000002           1         2           getuid
  1.09    0.000005           3         2           prctl
```

```
$ sudo strace -cS calls lsof -n
...

% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
 15.98    0.143027           3     48627           write
 18.87    0.168962           3     48288       289 stat
 24.97    0.223523           9     23936           read
  6.01    0.053806           2     21532      1021 close
 10.05    0.090014           4     20540        36 openat
 10.13    0.090710           4     20216        85 readlink
  5.16    0.046198           2     19569           fstat
  7.39    0.066138           4     18840           lstat
  1.14    0.010197           7      1360           getdents
  0.09    0.000805           3       320           rt_sigaction
  0.09    0.000846           3       318           alarm
  0.07    0.000647           4       155           brk
```

Some common used syscalls:

```text
$ whatis --section 2 read write openat close stat fstat lstat mmap munmap mprotect socket ioctl fcntl futex select connect bind access execve sendmsg recvmsg clone brk
read (2)             - read from a file descriptor
write (2)            - write to a file descriptor
openat (2)           - open and possibly create a file
close (2)            - close a file descriptor
stat (2)             - get file status
fstat (2)            - get file status
lstat (2)            - get file status
mmap (2)             - map or unmap files or devices into memory
munmap (2)           - map or unmap files or devices into memory
mprotect (2)         - set protection on a region of memory
socket (2)           - create an endpoint for communication
ioctl (2)            - control device
fcntl (2)            - manipulate file descriptor
futex (2)            - fast user-space locking
select (2)           - synchronous I/O multiplexing
connect (2)          - initiate a connection on a socket
bind (2)             - bind a name to a socket
access (2)           - check user's permissions for a file
execve (2)           - execute program
sendmsg (2)          - send a message on a socket
recvmsg (2)          - receive a message from a socket
clone (2)            - create a child process
brk (2)              - change data segment size
```

For details, run `man 2 SYSCALL`, e.g `man 2 select`.

### Other useful usage of strace
#### Attach to a running process
```
strace -p PID
```
#### Follow forks
```
strace -f -p PID
```

### Scenarios
Some scenarios that strace come to shine:

- Learning under-the-hood what a program does
- Troulbeshooting a permission issue when a process failed without helpful
  error messages.
- See why a program stucks (looping, or waiting for something).

## Wrap up
- Program uses computer resources via system calls to tell Linux kernel.
- `strace` traces system calls
- system calll are documented at `man 2 SYSCALL`
