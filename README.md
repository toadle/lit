# lit

<p>
  <img src="https://res.cloudinary.com/toadle/image/upload/v1666520354/lit_uy82pn.png" width="300" alt="The lit rocket">
</p>

_"launch in time"_ - An easily configurable quicklauncher for your terminal. Think of it like Alfred, Launchbar or Raycast for your terminal.
Use it as a replacement for Spotlight and make it trigger via your Cmd+Space.

Type away into the fuzzy prompt to

- Launch apps
- Find contacts
- Translate
- Calculate
- Convert currencies

All the stuff you desired is directly integrated from you sheel via a simple YAML-config.

## Demo

[![asciicast](https://asciinema.org/a/c6dbWieV2Pn8DgfIbgqZlByPE.svg)](https://asciinema.org/a/c6dbWieV2Pn8DgfIbgqZlByPE)

## Installation

Download a binary from the [releases][releases] page.

Or build it yourself (requires Go 1.18+):

```bash
brew tap toadle/lit
brew install lit
```
## Uninstall

```bash
brew uninstall lit
brew untap toadle/lit
```

[releases]: https://github.com/mrusme/canard/releases

## User Manual

Everything `lit` does is based on regular terminal commands. `lit` will re-execute them on every change of your input and you can then navigate up and down to execute actions on the result.

You can launch it simply by typing `lit` in your terminal of choice.

### Configuration

`lit` is configured through a file at `~/.lit.yml`. If it does not exist you need to create it.

Here you can configure two kinds of sources for `lit`:

**Calculators**
They are displayed above the input of your query. They should do a basic transformation of the input like calculations or translations. The result is then reduced to **one line** and displayed.

**Searches**
They are displayed below the input of your query. They are supposed to return **several lines** of results which are then fuzzy-filtered by the query. It is possible to read the results through a regular expression in order to make it more usable.

A basic `.lit.yml` looks like this:

```yaml
calculators:
  - command: "..."
    action: "..."
    label: "..."
  - ...

searches:
  - command: "..."
    format: "..."
    action: "..."
    labels:
      title: "..."
      description: "..."
  - ...
```

Here is what those values do:

**Calculators**

`calculators` may contain several configurations which each will result in a new line getting added to `lit`s interface above the query input.

`command` should contain a terminal command that gets re-excuted whenever the input changes. The input-value can be substituted with a `{input}`-placeholder. The result of the terminal command will be reduced to one line.

`action` should contain a terminal command that gets executed when you select the entry in the interface. The result of the `command` can be substituted via a `{data}`-placeholder. After selection `lit` will close.

`label` if you do not provide a label then `lit` will display the `command` next to it's result. Sometimes the `command`s can be very messy and verbose. Therefore if you provide a label `lit` will display it instead. The current input can be substituted via `{input}`.

**Searches**

`searches` may contain several configurations which are all executed when `lit` launches and will NOT be re-executed when the input changes (that will change in the future). All results will be combined and displayed under the query-input.

`command` should contain a terminal commands that gets executed ONCE when `lit` launches. Every line of the result will be one result in `lit`'s interface.

`format` should contain a regular expression with capture groups that helps make sense of the `command`'s result. The regexp will be applied to every line separately. The capture-groups will be available as `{...}`-substitutions in the `action` and `labels`.

`action` should contain a terminal command that gets executed when you select one of the results in the interface. The result of the `command` can be substituted via all `{...}`-placeholders that are available through the `format`-regexp's capture group. After selection `lit` will close.

`labels` here you can configure how the results from `command` get displayed. `lit` displays a `title` and `description` for each. For both substitutions via `{...}` are possible. Both entries can be left empty. `lit` will then default to looking for a `{title}` and `{description}` capture group in `format` in order to display something readable.

## Example

This repository contains a [`.lit.example.yml`](.lit.example.yml) that is contains a pretty great starting-point for you own configuration. This example-configuration contains the following features:

- Translation using `trans` from [translate-shell
  ](https://github.com/soimort/translate-shell)

- Calculation using the phenomenal `qalc` from [Qalculate!](https://qalculate.github.io)

- Quick-launching apps using `mdfind` the command-line Spotlight-client on the mac

- Searching contacts for your addressbook via `mdfind`.

All of the above is made for using it on a Mac. If somebody could provide a good Linux-starting-point I'd be happy to add it.

## Good to know

### Using pipes (`|`) in terminal commands

`lit` supports using `|` in all `command`s in order to chain several commands for more flexibility. It does so by implementing it's own output-redirection. Keep in mind that `|` normally is a feature of your shell of choice and not the commands themselves. Therefore don't get to many ideas using `>` or `>>` in your commands as `lit` currently does not support those.

### Dealing with spaces

Spaces in commands normally delimit arguments. But sometimes spaces a also contained in paths or filenames. Therefore `lit` will treat spaces differently depending on where they appear: If a space appears in the configuration e.g. within a `command` or `action`-configuration they will be treated as argument-delimiters. When spaces appear in the results of a command or within the result of a `{...}`-substitution they will be kept as a single argument. If you need to extract results to several arguments of an `action` please consider using `xargs`.

## License

[MIT](https://github.com/toadle/lit/raw/master/LICENSE)

Built using the great technologies provided by [Charm](https://github.com/charmbracelet)
