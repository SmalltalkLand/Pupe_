#!/bin/bash
x=$(mktemp /tmp/XXXXXX)
unsquashfs /app/bin/index.sfs $x
exec bwrap --unshare-user --uid 0 --gid 0 --bind $x / /sbin/init $@
