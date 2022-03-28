windres resources.rc -o resources.syso
go build -ldflags "-extld=gcc -extldflags=resources.syso -H=windowsgui -s -w" main.go