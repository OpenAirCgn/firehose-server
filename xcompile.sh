# source this script to build versions of:
#   - firehose server
#   - open air simulator
# for plattforms:
#   - osx, linux(amd64 and arm5/raspi), windows

VERSION=`git describe --tags`
DATE=`date +%Y%m%d`
LDFLAGS="-X main.version=${VERSION}_${DATE}"

for os in darwin linux windows; do
	GOOS=${os} GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o firehose.${VERSION}.${os} cmd/firehose/firehose.go 
	GOOS=${os} GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o oa_sim.${VERSION}.${os} cmd/simulator/openair.go
	if [ $os == "linux" ]; then 
GOOS=${os} GOARCH=arm GOARM=5 go build -ldflags "${LDFLAGS}" -o firehose.${VERSION}.raspi cmd/firehose/firehose.go
GOOS=${os} GOARCH=arm GOARM=5 go build -ldflags "${LDFLAGS}" -o oa_sim.${VERSION}.raspi cmd/simulator/openair.go
	fi
done

