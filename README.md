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
$ brain_dump --help
USAGE: brain_dump [-c] [-e] [-p PROFILE] INPUT

INPUT must be a list of arguments or - to get the input from the STDIN.
The input can also contains arguments using key=value format to pass variables for the template.
If no input is given, the output file will be opened in a text editor.

OPTIONS:
  -c    print the default configuration to STDOUT
  -e    open the editor after writing the file
  -p string
        name of the profile to use (default "default")
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

- Pass variables to the template file

To fill the templates, the user can either add variables in a per profile
configuration or pass them at the command line.

`$HOME/.config/brain_dump/templates/snips.tmpl`
````tmpl
{{ if .desc }}
{{ .desc }}
{{ end }}
```{{ if .lang }}{{ .lang }}{{ end }}
{{ trim .input }}
```
````

`$HOME/.config/brain_dump/brain_dump.json`
```jsonc
{
  // ...
  "snip": {
    // ...
    "template_vars": {
      "lang": "raw"
    }
  },
  // ...
}
```

```bash
$ brain_dump -p snip lang=bash desc="Hello World in bash" echo Hello, World
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
