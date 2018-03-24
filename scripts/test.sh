code=$(sh ./scripts/md5sum.sh "$@")
./client/client --code $code