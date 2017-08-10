SET CGO_ENABLED=0
SET GOARCH=amd64

SET GOOS=windows
go build  -o ./dic.exe ../.

SET GOOS=darwin
go build -o ./dic_mac ../.

SET GOOS=linux
go build -o ./dic ../.
