# Storenv

Storenv is a simple secret sharing solution. Configure your secrets and save the configuration with your code using the version control system. Store the secrets with this CLI. Now your colleagues can easily fetch the secrets.

## Installation

```bash
$ export STORENV_VERSION=0.0.0-alpha; OS=osx; ARCH=aarch64; export STORENV_TAR=storenv-$STORENV_VERSION-$OS-$ARCH.tar.gz
$ wget https://github.com/hanchchch/storenv/releases/download/$STORENV_VERSION/$STORENV_TAR
$ tar -xvf $STORENV_TAR
$ mv storenv /usr/local/bin/storenv
```

## Usage

### Configuration

Configure storenv with yaml format. It tries to find `storenv.yaml` by default. You can override it by `-c` flag on execution.

```yaml
# storenv.yaml
storage:
  s3: # only s3 is supported for now.
    bucket: my-credentials # this field will be replaced by user-scoped configuration in future, so the bucket name also can be secured.
    prefix: storenv
    region: ap-northeast-2 # optional. It reads local aws configuration by default.

# list paths of your secrets
secrets:
  - .env
  - secrets.txt
```

### Store & Load

After the configuration, you can store the secrets in configured storage with `store` command.

```bash
$ storenv store
store 2 secrets into S3 storage
.env                           -> s3://my-credentials/storenv/.env
secrets.txt                    -> s3://my-credentials/storenv/secrets.txt
```

To load these secrets on somewhere else, use `load` command.

```bash
$ storenv load
load 2 secrets from S3 storage
s3://my-credentials/storenv/.env                 -> .env
s3://my-credentials/storenv/secrets.txt          -> secrets.txt
```
