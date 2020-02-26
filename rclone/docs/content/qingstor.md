---
title: "QingStor"
description: "Rclone docs for QingStor Object Storage"
date: "2017-06-26"
---

<i class="fas fa-hdd"></i> QingStor
---------------------------------------

Paths are specified as `remote:bucket` (or `remote:` for the `lsd`
command.)  You may put subdirectories in too, eg `remote:bucket/path/to/dir`.

Here is an example of making an QingStor configuration.  First run

    rclone config

This will guide you through an interactive setup process.

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
Choose a number from below, or type in your own value
[snip]
XX / QingStor Object Storage
   \ "qingstor"
[snip]
Storage> qingstor
Get QingStor credentials from runtime. Only applies if access_key_id and secret_access_key is blank.
Choose a number from below, or type in your own value
 1 / Enter QingStor credentials in the next step
   \ "false"
 2 / Get QingStor credentials from the environment (env vars or IAM)
   \ "true"
env_auth> 1
QingStor Access Key ID - leave blank for anonymous access or runtime credentials.
access_key_id> access_key
QingStor Secret Access Key (password) - leave blank for anonymous access or runtime credentials.
secret_access_key> secret_key
Enter a endpoint URL to connection QingStor API.
Leave blank will use the default value "https://qingstor.com:443"
endpoint>
Zone connect to. Default is "pek3a".
Choose a number from below, or type in your own value
   / The Beijing (China) Three Zone
 1 | Needs location constraint pek3a.
   \ "pek3a"
   / The Shanghai (China) First Zone
 2 | Needs location constraint sh1a.
   \ "sh1a"
zone> 1
Number of connnection retry.
Leave blank will use the default value "3".
connection_retries>
Remote config
--------------------
[remote]
env_auth = false
access_key_id = access_key
secret_access_key = secret_key
endpoint =
zone = pek3a
connection_retries =
--------------------
y) Yes this is OK
e) Edit this remote
d) Delete this remote
y/e/d> y
```

This remote is called `remote` and can now be used like this

See all buckets

    rclone lsd remote:

Make a new bucket

    rclone mkdir remote:bucket

List the contents of a bucket

    rclone ls remote:bucket

Sync `/home/local/directory` to the remote bucket, deleting any excess
files in the bucket.

    rclone sync /home/local/directory remote:bucket

### --fast-list ###

This remote supports `--fast-list` which allows you to use fewer
transactions in exchange for more memory. See the [rclone
docs](/docs/#fast-list) for more details.

### Multipart uploads ###

rclone supports multipart uploads with QingStor which means that it can
upload files bigger than 5GB. Note that files uploaded with multipart
upload don't have an MD5SUM.

### Buckets and Zone ###

With QingStor you can list buckets (`rclone lsd`) using any zone,
but you can only access the content of a bucket from the zone it was
created in.  If you attempt to access a bucket from the wrong zone,
you will get an error, `incorrect zone, the bucket is not in 'XXX'
zone`.

### Authentication ###

There are two ways to supply `rclone` with a set of QingStor
credentials. In order of precedence:

 - Directly in the rclone configuration file (as configured by `rclone config`)
   - set `access_key_id` and `secret_access_key`
 - Runtime configuration:
   - set `env_auth` to `true` in the config file
   - Exporting the following environment variables before running `rclone`
     - Access Key ID: `QS_ACCESS_KEY_ID` or `QS_ACCESS_KEY`
     - Secret Access Key: `QS_SECRET_ACCESS_KEY` or `QS_SECRET_KEY`

### Restricted filename characters

The control characters 0x00-0x1F and / are replaced as in the [default
restricted characters set](/overview/#restricted-characters).  Note
that 0x7F is not replaced.

Invalid UTF-8 bytes will also be [replaced](/overview/#invalid-utf8),
as they can't be used in JSON strings.

<!--- autogenerated options start - DO NOT EDIT, instead edit fs.RegInfo in backend/qingstor/qingstor.go then run make backenddocs -->
### Standard Options

Here are the standard options specific to qingstor (QingCloud Object Storage).

#### --qingstor-env-auth

Get QingStor credentials from runtime. Only applies if access_key_id and secret_access_key is blank.

- Config:      env_auth
- Env Var:     RCLONE_QINGSTOR_ENV_AUTH
- Type:        bool
- Default:     false
- Examples:
    - "false"
        - Enter QingStor credentials in the next step
    - "true"
        - Get QingStor credentials from the environment (env vars or IAM)

#### --qingstor-access-key-id

QingStor Access Key ID
Leave blank for anonymous access or runtime credentials.

- Config:      access_key_id
- Env Var:     RCLONE_QINGSTOR_ACCESS_KEY_ID
- Type:        string
- Default:     ""

#### --qingstor-secret-access-key

QingStor Secret Access Key (password)
Leave blank for anonymous access or runtime credentials.

- Config:      secret_access_key
- Env Var:     RCLONE_QINGSTOR_SECRET_ACCESS_KEY
- Type:        string
- Default:     ""

#### --qingstor-endpoint

Enter a endpoint URL to connection QingStor API.
Leave blank will use the default value "https://qingstor.com:443"

- Config:      endpoint
- Env Var:     RCLONE_QINGSTOR_ENDPOINT
- Type:        string
- Default:     ""

#### --qingstor-zone

Zone to connect to.
Default is "pek3a".

- Config:      zone
- Env Var:     RCLONE_QINGSTOR_ZONE
- Type:        string
- Default:     ""
- Examples:
    - "pek3a"
        - The Beijing (China) Three Zone
        - Needs location constraint pek3a.
    - "sh1a"
        - The Shanghai (China) First Zone
        - Needs location constraint sh1a.
    - "gd2a"
        - The Guangdong (China) Second Zone
        - Needs location constraint gd2a.

### Advanced Options

Here are the advanced options specific to qingstor (QingCloud Object Storage).

#### --qingstor-connection-retries

Number of connection retries.

- Config:      connection_retries
- Env Var:     RCLONE_QINGSTOR_CONNECTION_RETRIES
- Type:        int
- Default:     3

#### --qingstor-upload-cutoff

Cutoff for switching to chunked upload

Any files larger than this will be uploaded in chunks of chunk_size.
The minimum is 0 and the maximum is 5GB.

- Config:      upload_cutoff
- Env Var:     RCLONE_QINGSTOR_UPLOAD_CUTOFF
- Type:        SizeSuffix
- Default:     200M

#### --qingstor-chunk-size

Chunk size to use for uploading.

When uploading files larger than upload_cutoff they will be uploaded
as multipart uploads using this chunk size.

Note that "--qingstor-upload-concurrency" chunks of this size are buffered
in memory per transfer.

If you are transferring large files over high speed links and you have
enough memory, then increasing this will speed up the transfers.

- Config:      chunk_size
- Env Var:     RCLONE_QINGSTOR_CHUNK_SIZE
- Type:        SizeSuffix
- Default:     4M

#### --qingstor-upload-concurrency

Concurrency for multipart uploads.

This is the number of chunks of the same file that are uploaded
concurrently.

NB if you set this to > 1 then the checksums of multpart uploads
become corrupted (the uploads themselves are not corrupted though).

If you are uploading small numbers of large file over high speed link
and these uploads do not fully utilize your bandwidth, then increasing
this may help to speed up the transfers.

- Config:      upload_concurrency
- Env Var:     RCLONE_QINGSTOR_UPLOAD_CONCURRENCY
- Type:        int
- Default:     1

<!--- autogenerated options stop -->