# brain\_dump

`brain_dump` allows to quickly write notes, todos, snippets, etc. in text files.

Tempate files ([Golang text/template](https://pkg.go.dev/text/template) + [sprig](https://masterminds.github.io/sprig/))
can be used to control and format what will be written in the files.

## Project status

This project is in a very early stage.

I use this tool since the beginning and got only one issue, which is fixed by now.

No new features are planned right now unless some users have suggestions.

## Usage

- Print help/usage message.

```bash
$ brain_dump -h
$ brain_dump --help
```

- Dump a note without any configuration or any template.

Notice: if no configuration is found, the tool will print a log message
explaining that the default configuration will be used.

```bash
$ brain_dump the quick brown fox jumps over the lazy dog.
```

- Dump a note without any configuration or any template and then edit the file.

```bash
$ brain_dump -e the quick brown fox jumps over the lazy dog.
```

## Configuration

Configuration and template samples can be found in `./samples/`.

Here is the default configuration ( `brain_dump -c` ) :

```json
{
  "default": {
    "file": "~/Documents/brain_dump.md",
    "write_mode": "append",
    "template_file": "",
    "template_vars": null,
    "formats": {
      "date": "2006-01-02",
      "datetime": "2006-01-02 15:04:05",
      "time": "15:04:05"
    },
    "keys_case": "snake_case",
    "editor": "$EDITOR",
    "editor_args": null
  }
}
```

Here is an example of bash configuration :

```bash
if command -v brain_dump &> /dev/null; then
  alias bd=brain_dump
  alias link='brain_dump -p link'
  alias note='brain_dump -p note'
  alias snip='brain_dump -p snip'
  alias todo='brain_dump -p todo'
  export BRAINDUMP_CONFIG_FILE="$HOME/.config/brain_dump/brain_dump.json"
  export BRAINDUMP_DEBUG=1
fi
```
