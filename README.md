# wrcrux

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]

Add horcruxes (extension writers) to your writer.

An utility that allows you to pipe one writer to multiple writers (e.g. stdout, a file and a TCP connection). Can perform buffered or unbuffered per message.

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

[godoc-image]: https://godoc.org/github.com/avrebarra/wrcrux?status.svg
[godoc-url]: https://godoc.org/github.com/avrebarra/wrcrux
[report-image]: https://goreportcard.com/badge/github.com/avrebarra/wrcrux
[report-url]: https://goreportcard.com/report/github.com/avrebarra/wrcrux
[tests-image]: https://cloud.drone.io/api/badges/avrebarra/wrcrux/status.svg
[tests-url]: https://cloud.drone.io/avrebarra/wrcrux
[coverage-image]: https://codecov.io/gh/avrebarra/wrcrux/graph/badge.svg
[coverage-url]: https://codecov.io/gh/avrebarra/wrcrux
[sponsor-image]: https://img.shields.io/badge/github-donate-green.svg