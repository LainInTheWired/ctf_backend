#!/bin/bash

# RPCサービスの起動
service rpcbind start

# NFSサービスの起動
service nfs-kernel-server start

# NFSサービスをフォアグラウンドで実行
tail -f /dev/null
