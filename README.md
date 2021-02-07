# adi2cbr

A script for converting an ADIF amateur radio general log file to a Cabrillo amateur radio contest
file format.

Right now this is just a quick hack set up specifically for Winter Field Day and for my info. I
might generalize it someday.

## Usage

The program is set up like a Unix utility, so it reads from `STDIN` and outputs the Cabrillo
to `STDOUT`. It prints notes to `STDERR`.

```shell
$ go build
$ <wfd.adi adi2cbr >k0swe.cbr
```
