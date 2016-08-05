# isosim

Iso Websim is a ISO8583 simulator built using golang (http://golang.org). It provides a simple interface to load ISO specifications. 
The specifications themselves are defined in text file (more information on developing your own specs can be found here - 
../specs/isoSpecs.spec. 

You can load an existing trace or create a message from scratch and send it to a ISO host (identified by host and a port). 

Note: Please note that this application has been tested to work only on chrome browser.

The main program is started by the file isosim/isosim.go

C:\go run isosim.go -specDefFile specs\isoSpecs.spec -htmlDir .\html

Usage:
---------------------------------------
C:>go run isosim.go -help
Usage of C:\isosim.exe:
  -debugEnabled
        true if debug logging should be enabled. (default true)
  -htmlDir string
        Directory that contains any HTML's and js/css files etc. (default ".")
  -httpPort int
        Http port to listen on. (default 8080)
  -specDefFile string
        The file containing the ISO spec definitions. (default "isoSpec.spec")
exit status 2

