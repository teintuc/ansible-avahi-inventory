ansible-avahi-discovery
======

An ansible dynamic inventory based an [avahi(mdns)](https://www.avahi.org/). I first made it to configure my raspberry pi devices but it would be interesting for iot too.

Right now, it "shoud" works. It has not being fully tested, etc ... The code is not really clean and some improvements has to be made.

### Prerequies

#### On the devices you want to manage
You need to have an avahi-daemon started and configured on each devices you want to reach.

It should be seen by the command `avahi-browser -arp`

#### On your computer

You must have `golang` installed on your computer. The project has been developped using `go version go1.19.5 linux/amd64`.

### Compilation 

Then it is just a question of `make build` and you will find the binary in the `bin` directory. It makes a binary for Linux and MacOS.

### Use it

Once your binary is compiled, use that binary as an inventory in ansible: `ansible-inventory -i bin/avahi-discovery-linux-amd64 --list`
