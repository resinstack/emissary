# Emissary

Emissary is a small utility for one-shot templating of files with
secrets.  Its primary use case is machines on AWS or other clouds
where it is desireable to fetch secret data on startup from some
secure storage solution.  The binary uses templates as described below
to write out the secret data to files.

## Templates

Emissary uses templates with YAML front matter to control where
secrets are written.  Templates will be loaded from either the current
working directory, or from the path referred to by
`EMISSARY_BASE_PATH`, and must satisfy the glob pattern `*.tpl`.

An example of a template is shown below:

```
---
mode: 0400
dest: /run/config/nomad/secret.hcl
onrender: service nomad restart
---
consul {
  token = "{{poll "insecure" "consul_token"}}"
}
```

Values contained in the frontend are as follows:

  * `mode`: File mode to write the rendered template as, expressed in 4
    digit octal notation.
  * `dest`: Path to which the rendered template will be written.  This
    path must be writable to the emissary user.
  * `onrender`: An optional command which will be executed after the
    template is written to disk.  This is intended to provide a way to
    signal or restart the application which will consume the secret
    data.

Secrets may be added to templates by using a function that summons
them.  Currently implemented secret summoners are:

  * `poll`: The poll function will attempt to get the secret value
    every 10 seconds and will retry on all failures except connection
    problems or problems with permissions.  Importantly, secrets that
    do not exist and secrets that contain empty values will be retried
    until they exist or contain non-empty values, respectively.

Loading a secret uses a template expression as follows:

```
{{ <method> "<system>" "<id>" }}
```

The method is a secret summoner as described above.  The system
identifies which secret storage system to use.  Currently implemented
secret storage systems are described below.  The ID field is the name
of the secret as will be recognized by the storage system.

## Secret Stores

Secret stores are pluggable components that provide a way to obtain a
secret from the remote service.  These components may themselves need
authorization, how this authorization is obtained is beyond the scope
of Emissary.

  * `insecure`: The insecure storage engine queries a remote webserver
    for secrets.  It makes requests of the form `$INSECURE_BASE/<id>`
    where `$INSECURE_BASE` is an environment variable containing the
    base URL of the server to query.
  * `awssm`: Fetches secrets from the AWS Secrets Manager.
    Credentials are loaded from the normal locations for the AWS SDK,
    and the region can be deduced from a running EC2 instance if
    applicable.
