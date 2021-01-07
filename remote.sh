#!/bin/sh
source rbuild.env
rsync --exclude=bin --exclude=vendor  -arv . root@${REMOTE_HOST}:${REMOTE_WD}/ --delete
ssh root@${REMOTE_HOST} ". /etc/profile; cd ${REMOTE_WD}/; make build"
scp root@${REMOTE_HOST}:${REMOTE_WD}/tmp/boot.iso "${ISO_DIR}/${ISO_NAME}"
