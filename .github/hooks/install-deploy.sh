#!/usr/bin/env bash
# deploy in linux, ci/cd

cd /tmp/dist/ || exit
TARGET_OS=`uname | tr '[:upper:]' '[:lower:]'`
TARGET_ARCH=`uname -m`
[ "${TARGET_ARCH}" == 'x86_64' ] && TARGET_ARCH=amd64
ls -t bypctl-*-${TARGET_OS}-${TARGET_ARCH}.tar.gz | tail -n +2 | xargs rm -f
tar -zxvf bypctl-*-${TARGET_OS}-${TARGET_ARCH}.tar.gz
cd bypctl-*-${TARGET_OS}-${TARGET_ARCH}
\mv bypctl /data/webroot/mirrors/bypanel/bypctl-${TARGET_OS}-${TARGET_ARCH}
cd /data/webroot/mirrors/bypanel
sed -i "s@.*bypctl-${TARGET_OS}-${TARGET_ARCH}@$(md5sum bypctl-${TARGET_OS}-${TARGET_ARCH})@" /data/webroot/mirrors/md5sum.txt
rm -rf /tmp/dist
