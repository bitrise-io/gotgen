# gotgen

Minimalistic Go Template based Generator

## How to use

Run:

```
gotgen init
```

to create a `gg.conf.json` GotGen configuration file in the current directory.

Then run:

```
gotgen generate
```

Reads `gg.conf.json` from the currend directory, parses all the `.gg` files in the current directory,
then runs [Go template](https://golang.org/pkg/text/template/) on all `.gg` file content,
with Inventory defined in the `gg.conf.json` exposed as Inventory for the Go Template,
then saves the generated files with the same name without `.gg` extension.
