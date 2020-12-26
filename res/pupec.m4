#!/bin/bash
changeword(`@\([_a-zA-Z0-9]*\)')

emit(){
echo "$1" > /var/main/tmp/pupe/bus1
}
dev=/$(head -c 30 /var/main/dev/random | base64 -w0)
proc=/$(head -c 30 /var/main/dev/random | base64 -w0)
exec bwrap --bind / / --bind /var/main/dev $dev --proc $proc --bind /var/main $dev/host /bin/pc2 $dev $proc $@
