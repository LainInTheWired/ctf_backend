# ベースイメージとして公式Ubuntuを使用
FROM ubuntu:22.04

# 環境変数の設定（非対話モード）
ENV DEBIAN_FRONTEND=noninteractive

# 必要なパッケージのインストール
RUN apt-get update && \
    apt-get install -y nfs-kernel-server rpcbind && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 共有ディレクトリの作成
RUN mkdir -p /nfsshare && \
    chown -R nobody:nogroup /nfsshare

# エクスポート設定ファイルの追加
COPY exports /etc/exports

# ポートのエクスポート
EXPOSE 2049 111 20048 32769

# スタートアップスクリプトの追加
COPY start-nfs.sh /start-nfs.sh
RUN chmod +x /start-nfs.sh

# コンテナ起動時に実行されるコマンド
CMD ["/start-nfs.sh"]
