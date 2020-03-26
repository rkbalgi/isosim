[![Go Report Card](https://goreportcard.com/badge/github.com/rkbalgi/isosim)](https://goreportcard.com/report/github.com/rkbalgi/isosim)
[![codecov](https://codecov.io/gh/rkbalgi/isosim/branch/master/graph/badge.svg)](https://codecov.io/gh/rkbalgi/isosim)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/rkbalgi/isosim?tab=doc)
![build-isosim-workflow](https://github.com/rkbalgi/isosim/workflows/build-isosim-workflow/badge.svg?branch=master)


# ISO WebSim

![](https://github.com/rkbalgi/isosim/blob/master/docs/images/home.png)


* A quick demo - https://github.com/rkbalgi/isosim/wiki/Test-Examples
* Running on Docker - https://github.com/rkbalgi/isosim/wiki/Running-on-Docker


Iso Websim is a ISO8583 simulator built using golang (http://golang.org). 

It provides 
* A simple mechanism to define ISO specifications 
* Define servers based on defined specs to respond to incoming messages (rules based on amount etc)
* Build and send transactions to servers (as a client)
* Save messages that be can be replayed later

The specifications themselves are defined in text file (more information on developing your own specs can be found in (https://github.com/rkbalgi/isosim/blob/master/specs/isoSpecs.spec).

A new front end is being developed in react. See [here](#new-development-of-a-frontend-in-React)

The main program is started by the file cmd/isosim/isosim.go
 
* __:warning: Please note that this application has been tested on chrome browser only.__

### Usage: 
```
C:>go run isosim.go -help
  -dataDir string
        Directory to store messages (data sets). This is a required field.
  -debugEnabled
        true if debug logging should be enabled. (default true)
  -htmlDir string
        Directory that contains any HTML's and js/css files etc. (default ".")
  -httpPort int
        Http port to listen on. (default 8080)
  -specDefFile string
        The file containing the ISO spec definitions. (default "isoSpec.spec")
exit status 2
```

### Starting Isosim 
```
$> set GOPATH=<Your Directory>
$> cd src\github.com\rkbalgi\isosim\cmd\isosim
$> go run isosim.go -httpPort 8080 -specDefFile ..\..\specs\isoSpecs.spec -htmlDir ..\..\html --dataDir ..\..\testdata
```
And now open chrome and hit this URL [http://localhost:8080/iso/home](http://localhost:8080/iso/home)

Read more about this on the [wiki](https://github.com/rkbalgi/isosim/wiki)

## New development of a frontend in React
The below is a screenshot of the revamped application with the frontend written in React.

![](https://github.com/rkbalgi/isosim/blob/master/docs/images/ReactApp_Screenshot.png)

The frontend is bundled with the application and can be accessed at [http://localhost:8080/]


