# autoeq-to-pulseeffects

A small Go application to convert [AutoEq][autoeq] TXT files and patch them
over PulseEffects' JSON files.

[autoeq]: https://github.com/jaakkopasanen/AutoEq

## Installation

```sh
go get github.com/diamondburned/autoeq-to-pulseeffects
```

## Usage

```sh
cd ~/.config/PulseEffects/output
# With patching and output file:
autoeq-to-pulseeffects -p "Old Preset.json" -o "New Preset.json" "~/Downloads/Headphones.txt"
# Without patching and output to stdout:
autoeq-to-pulseeffects "~/Downloads/Headphones.txt"
```
