# template

The `template` package provides a lightweight secret-injection renderer for
vaultpull. It replaces `{{KEY}}` placeholders inside any text file with
corresponding values fetched from Vault.

## Usage

```go
r := template.NewRenderer()
out, missing, err := r.RenderFile("config.tmpl", secrets)
```

## Placeholder syntax

Placeholders must use uppercase letters, digits, and underscores:

```
DATABASE_URL={{DB_HOST}}:{{DB_PORT}}/{{DB_NAME}}
API_KEY={{API_KEY}}
```

## CLI

```
vaultpull render config.tmpl --output config.env
```

If `--output` is omitted the rendered content is printed to stdout.

A warning is printed to stderr for every placeholder that could not be
resolved, but the command still exits successfully so pipelines are not broken.
