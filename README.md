# nm_dmenu
Simple programm for managing NetworkManager connections via dmenu written on golang
## Features
- Connect to exisiting NM wifi or wired connections
- Connect to new wifi connections that dont require passphrase
- Connect to _existing_ VPN connections 
- Enable/Disable networking
## License
- GPL3
## Requirements
1. Golang 1.18+
2. NetworkManager
3. Dmenu
## Installation
1. Make sure that you have golang installed, if it is go version should produce an output as showed below
```
$ go build
go version go1.20 linux/amd64
```
2. Download repository
```
git clone https://github.com/spanditime/nm_dmenu.git
```
3. Build 
```
go build
```
4. Copy nm_dmenu somwhere in your path
```
cp nm_dmenu /usr/bin/
```
5. Done :-)
## Usage
Simply launch nm_dmenu 
## Review
nm_dmenu was made to replace [twouters/nmcli-dmenu](https://github.com/twouters/nmcli-dmenu) because it wasn't behaving as i wanted in some scenarios:
1. Long launch times ( because of long wifi scan )
2. Clogged up list of wifi networks in public places

So i decided to make my own, on main page of menu there is only saved connections so it loads almost instantly, while wifi scans in background. Separate wifi page and a waiting page so you know that the rescan is actually in proccess, and you didn't missed your keybinding.
## TODO
- Add ability to delete connections
- Configure basic dmenu params through command line arguments( waiting title, font and color for different pages e.g. '-nf' change normfg for every page '-nfm' change normfg for main page only )
- Configure some dmenu patches params through command line arguments ( centered, xyw )
