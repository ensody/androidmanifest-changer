#!/usr/bin/env bash
set -euxo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

"$DIR/build-common.sh"

source "$DIR/utils.sh"

cat >> ~/.bashrc <<EOF
if [ -z "\$BUILD_VERSION_CHECK_DONE" ] && ! diff -q "$DIR" ".devcontainer/" > /dev/null; then
  echo -e "\e[1m\e[31mThis container is outdated. Please rebuild.\e[0m" > /dev/stderr
fi
export BUILD_VERSION_CHECK_DONE=true
EOF

# Make sure a few common tools are installed
apt-get update
apt-get upgrade -y
apt-get install -y --no-install-recommends curl gettext git gnupg less procps apt-utils locales bash-completion zip unzip

# Enable bash auto completion
cat >> ~/.bashrc <<EOF
source /etc/profile.d/bash_completion.sh
EOF

# This fixed the locale and is required for diff-so-fancy
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
locale-gen en_US.UTF-8

# Install diff-so-fancy
DIFF_SO_FANCY_VERSION="1.2.6"
download https://raw.githubusercontent.com/so-fancy/diff-so-fancy/v${DIFF_SO_FANCY_VERSION}/third_party/build_fatpack/diff-so-fancy \
  ed9de2669c789d1aba8456d0a7cf95adb326e220c99af4336405f21add8f0852 /usr/bin/diff-so-fancy
chmod a+x /usr/bin/diff-so-fancy

# Install protoc
PROTOC_VERSION="3.18.1"
download_zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
  220bd1704c73dbf4d0a91399a2ecf9d19938b5cd80c8a38839a023d8b87bb772 /usr/bin/ bin/protoc
chmod a+x /usr/bin/protoc

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
