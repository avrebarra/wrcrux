# wrcrux

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Sponsor][sponsor-image]][sponsor-url]

Writers horcrux. An utility that allows you to pipe one writer to multiple writers (e.g. stdout, a file and a TCP connection). Can be buffered or unbuffered.

All `Write` calls are queued into a channel that is read in a separate goroutine. Once the channel receives new data it writes the data to all registered outputs.

## Installation

```shell
go get -u github.com/avrebarra/wrcrux/
```

## Example

```go
wx := wrcrux.New(wrcrux.ConfigWux{})

wx.AddWriter(os.Stdout)
wx.AddWriter(os.Stdout) // or use other writer

wx.Write([]byte("data"))
wx.WriteRich(wrcrux.BImmediate, []byte("data sync"))
wx.WriteRich(wrcrux.BNormal, []byte("data unbuffered"))messages
```


## Style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.
