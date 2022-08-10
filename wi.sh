#!/usr/bin/env bash

mkdir -p "$HOME/.bottle"
cd "$HOME/.bottle"

ARCH=$(uname -m)
OS=$(uname -s)


if [ $ARCH = "x86_64" ]; then
    ARCH="amd64"
fi

# echo "Installing bottle-$OS-$ARCH"

install_latest() {
    echo -e "Installing bottle-$OS-$ARCH in \e[32m$HOME/.bottle/bottle\e[0m"
    if [ "$(uname)" = "Darwin" ]; then
    # Do something under Mac OS X platform 
        wget -q --show-progress --max-redirect 5 -O bottle "https://github.com/tristanisham/bottle/releases/latest/download/$1"
        chmod +x bottle
       
    elif [ $OS = "Linux" ]; then   
     # Do something under GNU/Linux platform
        wget -q --show-progress --max-redirect 5 -O bottle "https://github.com/tristanisham/bottle/releases/latest/download/$1"
        chmod +x bottle

    elif [ $OS = "MINGW32_NT" ]; then
    # Do something under 32 bits Windows NT platform
        curl "https://github.com/tristanisham/bottle/releases/latest/download/$1 -o bottle.exe"

    elif [ $OS == "MINGW64_NT" ]; then
    # Do something under 64 bits Windows NT platform
        curl "https://github.com/tristanisham/bottle/releases/latest/download/$1 -o bottle.exe"

    fi
}



if [ "$(uname)" = "Darwin" ]; then
    # Do something under Mac OS X platform 
    install_latest "bottle-darwin-$ARCH"       
elif [ $OS = "Linux" ]; then   
     # Do something under GNU/Linux platform
    install_latest "bottle-linux-$ARCH"
elif [ $OS = "MINGW32_NT" ]; then
    # Do something under 32 bits Windows NT platform
    install_latest "bottle-windows-$ARCH.exe"
elif [ $OS == "MINGW64_NT" ]; then
    # Do something under 64 bits Windows NT platform
    install_latest "bottle-windows-$ARCH.exe"
fi

echo
echo "Append the following to your $HOME/.profile or $HOME/.bash_rc"
echo
echo -e "export PATH=\$PATH:$HOME/.bottle"
echo