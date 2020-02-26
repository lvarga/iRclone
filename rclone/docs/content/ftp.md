---
title: "FTP"
description: "Rclone docs for FTP backend"
date: "2017-01-01"
---

<i class="fa fa-file"></i> FTP
------------------------------

FTP is the File Transfer Protocol. FTP support is provided using the
[github.com/jlaffaye/ftp](https://godoc.org/github.com/jlaffaye/ftp)
package.

Here is an example of making an FTP configuration.  First run

    rclone config

This will guide you through an interactive setup process. An FTP remote only
needs a host together with and a username and a password. With anonymous FTP
server, you will need to use `anonymous` as username and your email address as
the password.

```
No remotes found - make a new one
n) New remote
r) Rename remote
c) Copy remote
s) Set configuration password
q) Quit config
n/r/c/s/q> n
name> remote
Type of storage to configure.
Enter a string value. Press Enter for the default ("").
Choose a number from below, or type in your own value
[snip]
XX / FTP Connection
   \ "ftp"
[snip]
Storage> ftp
** See help for ftp backend at: https://rclone.org/ftp/ **

FTP host to connect to
Enter a string value. Press Enter for the default ("").
Choose a number from below, or type in your own value
 1 / Connect to ftp.example.com
   \ "ftp.example.com"
host> ftp.example.com
FTP username, leave blank for current username, ncw
Enter a string value. Press Enter for the default ("").
user> 
FTP port, leave blank to use default (21)
Enter a string value. Press Enter for the default ("").
port> 
FTP password
y) Yes type in my own password
g) Generate random password
y/g> y
Enter the password:
password:
Confirm the password:
password:
Use FTP over TLS (Implicit)
Enter a boolean value (true or false). Press Enter for the default ("false").
tls> 
Remote config
--------------------
[remote]
type = ftp
host = ftp.example.com
pass = *** ENCRYPTED ***
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

This remote is called `remote` and can now be used like this

See all directories in the home directory

    rclone lsd remote:

Make a new directory

    rclone mkdir remote:path/to/directory

List the contents of a directory

    rclone ls remote:path/to/directory

Sync `/home/local/directory` to the remote directory, deleting any
excess files in the directory.

    rclone sync /home/local/directory remote:directory

### Modified time ###

FTP does not support modified times.  Any times you see on the server
will be time of upload.

### Checksums ###

FTP does not support any checksums.

#### Restricted filename characters

In addition to the [default restricted characters set](/overview/#restricted-characters)
the following characters are also replaced:

File names can also not end with the following characters.
These only get replaced if they are last character in the name:

| Character | Value | Replacement |
| --------- |:-----:|:-----------:|
| SP        | 0x20  | ␠           |

Note that not all FTP servers can have all characters in file names, for example:

| FTP Server| Forbidden characters |
| --------- |:--------------------:|
| proftpd   | `*`                  |
| pureftpd  | `\ [ ]`              |

### Implicit TLS ###

FTP supports implicit FTP over TLS servers (FTPS). This has to be enabled
in the config for the remote. The default FTPS port is `990` so the
port will likely have to be explictly set in the config for the remote.

<!--- autogenerated options start - DO NOT EDIT, instead edit fs.RegInfo in backend/ftp/ftp.go then run make backenddocs -->
### Standard Options

Here are the standard options specific to ftp (FTP Connection).

#### --ftp-host

FTP host to connect to

- Config:      host
- Env Var:     RCLONE_FTP_HOST
- Type:        string
- Default:     ""
- Examples:
    - "ftp.example.com"
        - Connect to ftp.example.com

#### --ftp-user

FTP username, leave blank for current username, $USER

- Config:      user
- Env Var:     RCLONE_FTP_USER
- Type:        string
- Default:     ""

#### --ftp-port

FTP port, leave blank to use default (21)

- Config:      port
- Env Var:     RCLONE_FTP_PORT
- Type:        string
- Default:     ""

#### --ftp-pass

FTP password

- Config:      pass
- Env Var:     RCLONE_FTP_PASS
- Type:        string
- Default:     ""

#### --ftp-tls

Use FTP over TLS (Implicit)

- Config:      tls
- Env Var:     RCLONE_FTP_TLS
- Type:        bool
- Default:     false

### Advanced Options

Here are the advanced options specific to ftp (FTP Connection).

#### --ftp-concurrency

Maximum number of FTP simultaneous connections, 0 for unlimited

- Config:      concurrency
- Env Var:     RCLONE_FTP_CONCURRENCY
- Type:        int
- Default:     0

#### --ftp-no-check-certificate

Do not verify the TLS certificate of the server

- Config:      no_check_certificate
- Env Var:     RCLONE_FTP_NO_CHECK_CERTIFICATE
- Type:        bool
- Default:     false

#### --ftp-disable-epsv

Disable using EPSV even if server advertises support

- Config:      disable_epsv
- Env Var:     RCLONE_FTP_DISABLE_EPSV
- Type:        bool
- Default:     false

<!--- autogenerated options stop -->

### Limitations ###

Note that since FTP isn't HTTP based the following flags don't work
with it: `--dump-headers`, `--dump-bodies`, `--dump-auth`

Note that `--timeout` isn't supported (but `--contimeout` is).

Note that `--bind` isn't supported.

FTP could support server side move but doesn't yet.

Note that the ftp backend does not support the `ftp_proxy` environment
variable yet.

Note that while implicit FTP over TLS is supported,
explicit FTP over TLS is not.