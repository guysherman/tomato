# Tomato

An app for timing periods of focused work, on the terminal.

## Build

```
go get
go build .
```

## Usage

The main thing you can do via commandline arguments is specify durations for the focused work periods, 
as well as the short and long break periods. They take the form `<number><unit>`, where unit is one of: 
`ns`, `us`, `ms`, `s`, `m`, `h`. You can also do combinations.

Some examples:
* `300ms`
* `2h45m`

Tomato takes the following commandline args:

* `-f` the duration for the focused work periods
* `-s` the duration for the short break
* `-l` the duration for the long break
* `-L` the number of tomatos required to earn a long break

Focus Mode:
![A screenshot of Focus Mode](/doc/FocusMode.png)

Break Mode:
![A screenshot of Break Mode](/doc/BreakMode.png)


