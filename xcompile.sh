# source this script to build versions of:
#   - firehose server
#   - open air simulator
# for plattforms:
#   - osx, linux(amd64 and arm5/raspi), windows

VERSION=`git describe --tags`
DATE=`date +%Y%m%d`
LDFLAGS="-X main.version=${VERSION}_${DATE}"

REL_DIR=release
if [ ! -d ${REL_DIR} ]; then
	mkdir ${REL_DIR}
fi

for os in darwin linux windows; do
	GOOS=${os} GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${REL_DIR}/firehose.${VERSION}.${os} cmd/firehose/firehose.go 
	GOOS=${os} GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${REL_DIR}/oa_sim.${VERSION}.${os} cmd/simulator/openair.go
done

