PATH=/sbin:/usr/sbin:/bin:$PATH
m4  -I ./m4/m4/examples/ ./build.ninja.m4 > ./build.ninja
bwrap --unshare-user --uid 0 --gid 0 --bind / / ninja ./files
