# Go Xbox Controller Recording Passthrough
Use a Raspberry Pi Zero to connect to a wireless xbox controller, record input, and play back input by connecting raspberry pi as an emulated controller. Written in golang.

Use cases:
* Automated rune farming in Elden Ring

## Setup 
1. Install `raspberrypi-kernel-headers`
```bash
sudo apt install -y raspberrypi-kernel-headers
```
2. Get [xpadneo](https://github.com/atar-axis/xpadneo) xbox controller driver
```bash
# clone github.com/atar-axis/xpadneo repo
git clone https://github.com/atar-axis/xpadneo.git
```
3. Install driver
```bash
cd xpadneo

./install.sh
```
4. Pair controller over bluetooth with the Desktop interface or use the following commands if runing headless. You only need to do this step once if the controller is not already paired.
```bash
# launch bluetootctl
bluetoothctl

# enable scanning
scan on

# find your device MAC address - "Xbox Wireless Controller" and then issue pair command
pair <MAC Address>

# trust device
trust <MAC Address>

# connect device
connect <MAC Address>

# exit
exit
```
5. Test to ensure proper setup
```bash
go install github.com/0xcafed00d/joystick/joysticktest@master

$HOME/go/bin/joysticktest
```

## Run 
```bash
go run cmd/emulate/emulate.go
```