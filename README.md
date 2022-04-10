# brain\_dump

`brain_dump` allows to quickly write notes, todos, snippets, etc. in text files
using Golang [text/template](https://pkg.go.dev/text/template) and functions from the [Sprig](https://masterminds.github.io/sprig/) project.

## Project status

This project is in a very early stage. Things might break !

## Usage

- Print help/usage message.

```bash
$ brain_dump -h
$ brain_dump --help
```

- Print the default configuration to `stdout`.

```bash
$ brain_dump -c
```

- Dump a note without any configuration.

Notice: if no configuration is found, the tool will print a log message
explaining that the default configuration will be used.

```bash
$ brain_dump the quick brown fox jumps over the lazy dog.
```

- Dump a note without any configuration and then edit the file.

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
  export BRAINDUMP_CONFIG_HOME="$HOME/.config/brain_dump"
  export BRAINDUMP_CONFIG_FILE="$BRAINDUMP_CONFIG_HOME/brain_dump.json"
  export BRAINDUMP_DATA_HOME="$HOME/Documents/brain_dump"
  export BRAINDUMP_DEBUG=1
fi
```
