REM "---- Running ISO Websim .. ----"
set TLS_ENABLED=false
set TLS_CERT_FILE=C:\Users\rkbal\IdeaProjects\isosim-prj\src\isosim\certs\cert.pem
set TLS_KEY_FILE=C:\Users\rkbal\IdeaProjects\isosim-prj\src\isosim\certs\key.pem
go run isosim.go -http-port 8080 -specs-dir ..\..\specs -html-dir ..\..\html --data-dir ..\..\testdata
