REM "---- Starting ISO Websim .. ----"
REM "---- Setting ENV variables -----"
set TLS_ENABLED=false
set TLS_CERT_FILE=C:\Users\rkbal\IdeaProjects\isosim-prj\src\isosim\certs\cert.pem
set TLS_KEY_FILE=C:\Users\rkbal\IdeaProjects\isosim-prj\src\isosim\certs\key.pem
REM "---- Starting App -----"
go run isosim.go -http-port 8080 --log-level TRACE -specs-dir ..\..\test\testdata\specs -html-dir ..\..\web -data-dir ..\..\test\testdata\appdata
