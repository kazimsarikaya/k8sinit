#!/bin/sh -eu

cmd=${1:-build}

if [ "x$cmd" == "xbuild" ]; then
  REV=$(git describe --long --tags --match='v*' --dirty 2>/dev/null || git rev-list -n1 HEAD)
  NOW=$(date +'%Y-%m-%d_%T')
  GOV=$(go version)
  go mod tidy
  go mod vendor
  export CGO_ENABLED=0
  LDFLAGS="-w -s -extldflags=-static"
  for proto in $(find internal -name *.proto); do
    protoc --experimental_allow_proto3_optional -I $(dirname $proto) --go_out=$(dirname $proto) $(basename $proto)
  done
  go build -ldflags "${LDFLAGS} -X main.version=$REV -X main.buildTime=$NOW -X 'main.goVersion=${GOV}'"  -o ./bin/init ./cmd

  for m in $(lsmod |awk '{print $1}'|grep -v Module); do find /lib/modules/5.4.84-0-lts/ -name "$m.ko"; done |sort|sed -r "s%^/lib/modules/$(uname -r)/%%g" > hack/mkinitfs/features.d/k8sinit.modules
  rm -fr tmp/*
  mkdir -p tmp/iso/syslinux
  mkinitfs -o tmp/initramfs -P `pwd`/hack/mkinitfs/features.d/ -c `pwd`/hack/mkinitfs/mkinitfs.conf  -i `pwd`/bin/init
  cp -arv tmp/initramfs tmp/iso/
  cp -arv /boot/vmlinuz-lts tmp/iso/vmlinuz
  cp -arv /usr/share/syslinux/isolinux.bin tmp/iso/syslinux/
  cp -arv /usr/share/syslinux/ldlinux.c32 tmp/iso/syslinux/
  cp -arv hack/syslinux.cfg tmp/iso/syslinux/
  CURDIR=`pwd`
  cd tmp/iso
  genisoimage -J -l -o ../boot.iso -b syslinux/isolinux.bin -c syslinux/isolinux.cat -no-emul-boot -boot-load-size 4 -boot-info-table .
  cd $CURDIR
elif [ "x$cmd" == "xtest" ]; then
  shift
  ./test.sh $@
else
  echo unknown command
fi
