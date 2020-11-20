# wrcrux

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Sponsor][sponsor-image]][sponsor-url]

Add horcruxes to your writer. An utility that allows you to pipe multiple writers (e.g. stdout, a file and a TCP connection) from one writer. Can perform buffered or unbuffered per message.

All `Write` calls are queued into a channel that is read in a separate goroutine. Once the channel receives new data it writes the data to all registered outputs.

## Installation

```shell
go get -u github.com/avrebarra/wrcrux/
```

## Example

```go
wx := wrcrux.New(wrcrux.ConfigWux{})

wx.AddWriter(os.Stdout)
wx.AddWriter(os.Stdout)
wx.AddWriter(os.Stdout)

wx.Write([]byte("data"))
wx.WriteRich(wrcrux.Immediate, []byte("data sync"))
wx.WriteRich(wrcrux.Buffered, []byte("data unbuffered"))
```
