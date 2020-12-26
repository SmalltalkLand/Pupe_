#!/bin/ash
changeword(`@\([_a-zA-Z0-9]*\)')
@define(`console',`exec > /dev/console
exec 2>/dev/console
exec </dev/console')
@define(`shx',`-c "x=$(mktemp /tmp/XXXXXX);echo \"set -- $@;\" > $x;cat $x /proc/self/fd/0 | exec;rm -f $x;exec /proc/self/exe"')
PATH=/bin:/usr/bin:/sbin:/usr/sbin
if [ "$UID" != "0"];
if bwrap --bind / / true;
exec bwrap --bind / / --unshare-user --uid 0 --gid 0 /sbin/init $@
fi;
fi;
exec /bin/bash @shx $@ <<"ENE"
@console
y="/$(head -c 30 /dev/random | base64 -w0)"
mkfifo $y
mkfifo /bin/sock/reboot
/bin/shsrvr server /bin/sock/reboot /bin/reboot &
mkfifo /bin/sock/curl
/bin/shsrvr server /bin/sock/curl /usr/bin/curl &
mkfifo /bin/sock/pspy
/bin/shsrvr server /bin/sock/pspy /bin/pspy &
/bin/shsrvr server $y /bin/bash &
exec 1>(tmux new -s core /bin/shsrvr client $y @shx $@)
exec 0<&1
exec 3<&0
x(){
exec /bin/cat /dev/stdin <<"ENC"
exec 3</proc/1/fd/3
exec /bin/bash @shx $@ <<"ENB"
exec 0<&3
z=""
while ["$z" != "done"];
mount -o loop /dev/$z /mnt/pupe/$z
z="$( (echo "choose devices" && find /dev && echo "done") | sentaku)"
end;
w=""
w="$( (echo "choose save file:" && ls /dev) | sentaku)"
mount $w /save
x=""
y="/"
while ["$x" != "done"];
y="$y;$x"
a="$( (echo "choose a system file" && find -type f -name *.sfs /sfs /mnt/pupe /save/sfs /save/mnt/pupe && echo "done") | sentaku)"
if ["$a" != "done"]
x="/$(head -c 30 /dev/random | base64 -w0)"
mount -t squashfs $a $x
else
x="done"
fi
end;
mount -t overlay lowerdir=$y;/,upperdir=/save,workdir=/save/work /x
mount --bind /dev /x/dev
mount --bind /proc /x/proc
mount --bind /sys /x/sys
mount --bind /tmp/x /x/tmp
exec chroot /x /bin/bash @shx $@ <<"END"
exec 0<&3
mount -t overlay lowerdir=$(IFS=";" ls /var/pmod | xargs -i /bin/bash -c \'if [ "$(echo -e "load\$1\\\nyes\\\nno" | sentaku)" == "yes"];echo $1;else;echo "/";fi;\' {});/,upperdir=/,workdir=/_work /y
exec chroot /y /bin/bash @shx $@ <<"EN_D"
exec 0<&3
exec login root <<ENF
exec /bin/base64 -d | /bin/bash @shx $@ <<ENH$(cat <<"ENG" | base64)
ENH
exec 0<&3
exec unshare -m /bin/bash @shx $@ <<"ENI"
exec 0<&3
bind_port=$RANDOM
mount -t overlay lowerdir=/var/c/main;/,upperdir=/tmp/cmain,workdir=/tmp/cmain.work /var/c_/main
mount --bind /dev /var/c_/main/dev
mount --bind /proc /var/c_/main/proc
mount --bind /sys /var/c_/main/sys
bwrap --bind /var/c_/main / --bind / /var/main pupec docker_ build -t pupe/node /var/main/etc/pupe-node/Dockerfile
mkfifo /bin/sock/rofi.sock
bwrap --bind /var/c_/main / --bind / /var/main pupec node priveleged /bin/bash -c "eval \$(echo $(cat <<"ENJ" | base64 -w0) | base64 -d)"
openbox &
ENJ
bwrap --bind /var/c_/main / --bind / /var/main pupec node priveleged /bin/bash -c "eval \$(echo $(cat <<"ENJ" | base64 -w0) | base64 -d)"
apt-get install rofi
/var/host/bin/shsrvr server /var/host/bin/sock/rofi.sock rofi -dmenu &
ENJ
cat <<"ENJ" > /bin/rmenu
#!/bin/ash
exec bwrap --bind /var/c_/main / --bind / /var/main pupec node $(ls /container | /bin/shsrvr /bin/sock/rofi.sock) /bin/r $@
ENJ
mkfifo /bin/sock/flatpak.sock
bwrap --bind /var/c_/main / --bind / /var/main pupec node priveleged /bin/bash -c "eval \$(echo $(cat <<"ENJ" | base64 -w0) | base64 -d)"
apt-get install flatpak bubblewrap
/var/host/bin/shsrvr server /var/host/bin/sock/flatpak.sock bwrap --bind / / --bind /var/host/var/lib/flatpak /var/lib/flatpak flatpak &
ENJ
cat <<"ENJ" > /bin/flatpak
#!/bin/ash
exec /bin/shsrvr client /bin/sock/flatpak.sock $@
ENJ
/bin/busd &
exec /bin/binder insert 4 /bin/busybox $bind_port /bin/binder -- bwrap --bind /var/c_/main / --bind / /var/main /var/main/bin/binder recv 4 /bin/busybox -- 4 /bin/busybox ash -c "exec /var/main/bin/binder recv $bind_port /bin/binder -- 4 ash -c \$@" "exec bwrap --bind / / --bind /var/main/etc/init.d /etc/init.d --bind / /var/main --bind /var/main/tmp/.X11-unix /tmp/.X11-unix /bin/busybox \$@" init $@
ENI
ENG
ENF
EN_D
END
ENB
ENC
}
exec x
ENE
