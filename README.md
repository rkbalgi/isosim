[![Go Report Card](https://goreportcard.com/badge/github.com/rkbalgi/isosim)](https://goreportcard.com/report/github.com/rkbalgi/isosim)
[![codecov](https://codecov.io/gh/rkbalgi/isosim/branch/master/graph/badge.svg)](https://codecov.io/gh/rkbalgi/isosim)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/rkbalgi/isosim?tab=doc)
![build](https://github.com/rkbalgi/isosim/workflows/build/badge.svg)


# ISO WebSim


![](https://github.com/rkbalgi/isosim/blob/master/docs/images/home_rel2020.04.png)


Iso Websim is a ISO8583 simulator built using [golang](http://golang.org), [React](https://reactjs.org/), [material-ui](https://material-ui.com/) and
other amazing open source libraries.

Features -
* A mechanism to define ISO specifications
  * ASCII, EBCDIC, BCD and BINARY encoding for fields
  * Fixed, Variable, Bitmapped fields
  * Embedded/Nested fields 
  * Supported MLI's - 2I, 2E
* Define and run servers based on specs
  * Run servers from the UI or in [standalone mode](https://github.com/rkbalgi/isosim/wiki/Start-standalone-ISO-server-from-command-line)
  * Rules to respond to messages based on fields (rules based on amount, currency etc)   
* A UI to build and send transactions to servers (as a client)
  * Ability to edit fields on UI
  * Client-side validation of fields for content, length (more on the way) 
  * Save messages that be can be replayed later
* TLS, Docker 

The specifications themselves are defined in yaml file (Check out an example - [iso_specs.yaml](https://github.com/rkbalgi/isosim/blob/master/specs/iso_specs.yaml))

The [frontend](https://github.com/rkbalgi/isosim-react-frontend) is bundled with the application and can be accessed at [http://localhost:8080/](http://localhost:8080/)


* A quick demo - https://github.com/rkbalgi/isosim/wiki/Test-Examples
* Running on Docker - https://github.com/rkbalgi/isosim/wiki/Running-on-Docker

 
` Please note that this application has been tested on chrome browser only.`

### Usage: 
```
C:>go run isosim.go -help
  -data-dir string
        Directory to store messages (data sets). This is a required field.
  -debug-enabled
        true if debug logging should be enabled. (default true)
  -html-dir string
        Directory that contains any HTML's and js/css files etc.
  -http-port int
        Http port to listen on. (default 8080)
  -specs-dir string
        The directory containing the ISO spec definition files.
```

### Running Iso WebSim 
```
$> git checkout https://github.com/rkbalgi/isosim.git
$> cd isosim\cmd\isosim
$> go run isosim.go -httpPort 8080 -specs-dir ..\..\specs -html-dir ..\..\html -data-dir ..\..\testdata
```
Open chrome and hit this URL [http://localhost:8080/](http://localhost:8080/)

Read more about this on the [wiki](https://github.com/rkbalgi/isosim/wiki)

The old front end is still available at [http://localhost:8080/iso/home](http://localhost:8080/iso/home)



