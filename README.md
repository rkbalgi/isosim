[![Go Report Card](https://goreportcard.com/badge/github.com/rkbalgi/isosim)](https://goreportcard.com/report/github.com/rkbalgi/isosim)
[![codecov](https://codecov.io/gh/rkbalgi/isosim/branch/master/graph/badge.svg)](https://codecov.io/gh/rkbalgi/isosim)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/rkbalgi/isosim?tab=doc)
![build](https://github.com/rkbalgi/isosim/workflows/build/badge.svg)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/rkbalgi/isosim?include_prereleases&style=flat)
![Docker Pulls](https://img.shields.io/docker/pulls/rkbalgi/isosim?color=%23FF6528&label=docker%20pulls)

# ISO WebSim
A very short screencast - [https://youtu.be/vSRZ_nzU-Jg](https://youtu.be/vSRZ_nzU-Jg)

![](https://github.com/rkbalgi/isosim/blob/master/docs/images/home_rel2020.04_01.png)
![](https://github.com/rkbalgi/isosim/blob/master/docs/images/home_rel2020.04_02.png)


Iso Websim is a ISO8583 simulator built using [Go](http://golang.org), [React](https://reactjs.org/), [Material-UI](https://material-ui.com/) and
other amazing open source libraries.

## Features -
* A mechanism to define ISO specifications
  * ASCII, EBCDIC, BCD and BINARY encoding for fields
  * Fixed, Variable, Bitmapped fields
  * Embedded/Nested fields 
  * Supported MLI's - 2I, 2E, 4I, 4E
* Define and run servers based on specs
  * Run servers from the UI or in [standalone mode](https://github.com/rkbalgi/isosim/wiki/Start-standalone-ISO-server-from-command-line)
  * Rules to respond to messages based on fields (rules based on amount, currency etc)   
* A UI to build and send transactions to servers (as a client)
  * Ability to edit fields on UI
  * Client-side validation of fields for content, length (more on the way)
  * [PIN](https://github.com/rkbalgi/isosim/wiki/Pin-Field-(DF52)-Generation) and [MAC](https://github.com/rkbalgi/isosim/wiki/MAC-Generation-(DF64-128)) generation 
  * [Padding](https://github.com/rkbalgi/isosim/wiki/Field-Padding) support
  * Save messages that be can be replayed later
  * [Log](https://github.com/rkbalgi/isosim/wiki/Message-History) of past messages
* TLS, [Docker](https://github.com/rkbalgi/isosim/wiki/Running-on-Docker) 



Checkout the [wiki](https://github.com/rkbalgi/isosim/wiki) for more details!

The specifications themselves are defined in yaml file (Check out an example - [iso_specs.yaml](https://github.com/rkbalgi/isosim/blob/master/test/testdata/specs/iso_specs.yaml))

The [frontend](https://github.com/rkbalgi/isosim-react-frontend) is bundled with the application and can be accessed at [http://localhost:8080/](http://localhost:8080/)


* A quick walkthrough - https://github.com/rkbalgi/isosim/wiki/Test-Examples


 
` Please note that this application has been tested only on the chrome browser.`

### Usage: 
```
C:>go run isosim.go -help
  -data-dir string
        Directory to store messages (data sets). This is a required field.
  -html-dir string
        Directory that contains any HTML's and js/css files etc.
  -http-port int
        HTTP/s port to listen on. (default 8080)
  -log-level string
        Log level - [trace|debug|warn|info|error]. (default "debug")
  -specs-dir string
        The directory containing the ISO spec definition files.
```

### Running Iso WebSim 
```
$> git checkout https://github.com/rkbalgi/isosim.git
$> cd isosim\cmd\isosim
$> go run isosim.go -http-port 8080 -specs-dir ..\..\test\testdata\specs -html-dir ..\..\web -data-dir ..\..\test\testdata\appdata
```
Open chrome and hit this URL [http://localhost:8080/](http://localhost:8080/)




