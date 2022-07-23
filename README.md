# Ping-TUI

## Run

```
go get
go build
sudo chown root ping-tui
sudo chmod u+s ping-tui
./ping-tui
```

## Usage

![Demo](https://imgur.com/RezqulK.gif)

### `Enter` in input window

Start ping

### `Ctrl+D`, `Ctrl+S` in input and output window

Stop ping

### `Tab`, `Shift+Tab` in input and output window

Switch between input and output window

### `Esc`, `Ctrl+C`, `Ctrl+Q` in input and output window

Quit the program

### `Ctrl+L` in output window

Clear output window

### Mouse Left Click

Select window

### Mouse Scroll in output window

Scroll to view the output

## Reference

- [ping example](https://golangexample.com/a-very-simple-and-small-ping-tool-that-sends-icmp-echo-datagram-to-a-host/)
- [ping code](https://github.com/z1cheng/c-ping/blob/master/ping.go)
- [checksum](https://inc0x0.com/icmp-ip-packets-ping-manually-create-and-send-icmp-ip-packets/)
- [wikipedia ping packet](https://en.wikipedia.org/wiki/Ping_(networking_utility)#ICMP_packet)
