# Gatekeeper
Welcome to the *Gatekeeper* module istallation guide. This package is written in Go programming language. It is a simple, lightweight authentication, authorization tool.

## Installation (Manual)
### Step 1: Install Go runtime
Please follow the download and install instrictions located here: https://golang.org/doc/install. Also you need to prepare your environment for optimal Go experience. Please read and follow the helpful instructions found here: https://golang.org/doc/code.html

### Step 2: Installing Gatekeeper
After you have cloned the source code in the proper source directory structure as specified in the guides mentioned above, installing *Gatekeeper* is very easy. Simply follow these steps -

1. From inside the source folder where you copied the *Gatekeeper* files, simply run `go get`
2. Then run `go install`
3. If you followed the environment setup instructions, you should be able to launch *Gatekeeper* by simply typing `auth-utils`from any place.
4. Alternatively, from within the source folder where *Gatekeeper* code files were copied into, do `go run *.go` to launch the m-service.

```
The service will start at port 8000 if you did not change in the configuration file.
```

## Installation (Automated)
```
The installation script included has been tested on Ubuntu 14.04-03 64 bit release.
```

1. Make the script executable: chmod +x install.sh
2. Run it, sit back and relax: ./install.sh

## API and Usage Guide
Please see the T-Nova internal wiki page for API example snippets, [T-Nova Gatekeeper](http://wiki.t-nova.eu/tnovawiki/index.php/Gatekeeper)
Public API document will be made available very soon.

## Word of caution
This is a v0 release, the code is under active development with features being added rapidly.

```
Use gatekeeper.cfg to control program parameters such as ports, log and database files, etc.
```

## Credits
### Development Team
* Piyush Harsh (harh@zhaw.ch) / ICCLab