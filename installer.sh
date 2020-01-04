#!/bin/bash

# Based on a script found in https://github.com/rust-lang/rustup
set -u

GENESIS_URL=https://assets.infra.whiteblock.io/cli/bin/genesis

main() {
    downloader --check
    need_cmd uname
    need_cmd chmod
    need_cmd mkdir

    get_architecture || return 1
    local _arch="$RETVAL"
    assert_nz "$_arch" "arch"

    local _ext=""
    case "$_arch" in
        *windows*)
            _ext=".exe"
            ;;
    esac

    local _url="${GENESIS_URL}/${_arch}/genesis${_ext}"
    local _dir=$HOME/.whiteblock
    local _file="${_dir}/bin/genesis${_ext}"

    ensure mkdir -p "${_dir}/bin"
    ensure downloader "$_url" "$_file"
    ensure chmod a+x "$_file"
    if [ ! -x "$_file" ]; then
        printf 'Failed to install genesis\n'1>&2
        exit 1
    fi

    ignore "$_file" "$@"

    if [ ! -x "${_dir}/env" ]; then
        echo "export PATH=\"${_dir}/bin:$PATH\"" >> ${_dir}/env
        echo "source ${_dir}/env" >> $HOME/.bashrc

        echo "Please restart your shell or run"
        echo "source ${_dir}/env"
        echo "to start using genesis"
    fi

    

    local _retval=$?

    return "$_retval"
}

get_architecture() {
    local _ostype _cputype _bitness _arch
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"

    if [ "$_ostype" = Darwin ] && [ "$_cputype" = i386 ]; then
        # Darwin `uname -m` lies
        if sysctl hw.optional.x86_64 | grep -q ': 1'; then
            _cputype=x86_64
        fi
    fi

    case "$_ostype" in

        Linux)
            _ostype=linux
            ;;

        FreeBSD)
            _ostype=freebsd
            ;;

        NetBSD)
            _ostype=netbsd
            ;;

        DragonFly)
            _ostype=dragonfly
            ;;

        Darwin)
            _ostype=darwin
            ;;

        MINGW* | MSYS* | CYGWIN*)
            _ostype=windows-gnu
            ;;

        *)
            err "unrecognized OS type: $_ostype"
            ;;

    esac

    case "$_cputype" in

        i386 | i486 | i686 | i786 | x86)
            _cputype=386
            ;;

        xscale | arm | armv6l)
            _cputype=arm
            ;;

        armv7l | armv8l)
            _cputype=arm64
            ;;

        aarch64)
            _cputype=aarch64
            ;;

        x86_64 | x86-64 | x64 | amd64)
            _cputype=amd64
            ;;

        ppc)
            _cputype=ppc
            ;;

        ppc64)
            _cputype=ppc64
            ;;

        ppc64le)
            _cputype=ppc64le
            ;;

        s390x)
            _cputype=s390x
            ;;

        *)
            err "unknown CPU type: $_cputype"

    esac

    _arch="${_ostype}/${_cputype}"

    RETVAL="$_arch"
}

say() {
    printf 'genesis: %s\n' "$1"
}

err() {
    say "$1" >&2
    exit 1
}

need_cmd() {
    if ! check_cmd "$1"; then
        err "need '$1' (command not found)"
    fi
}

check_cmd() {
    command -v "$1" > /dev/null 2>&1
}

assert_nz() {
    if [ -z "$1" ]; then err "assert_nz $2"; fi
}

ensure() {
    if ! "$@"; then err "command failed: $*"; fi
}

ignore() {
    "$@"
}

downloader() {
    local _dld
    if check_cmd curl; then
        _dld=curl
    elif check_cmd wget; then
        _dld=wget
    else
        _dld='curl or wget' # to be used in error message of need_cmd
    fi

    if [ "$1" = --check ]; then
        need_cmd "$_dld"
    elif [ "$_dld" = curl ]; then
        if ! check_help_for curl --proto --tlsv1.2; then
            echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            curl --silent --show-error --fail --location "$1" --output "$2"
        else
            curl --proto '=https' --tlsv1.2 --silent --show-error --fail --location "$1" --output "$2"
        fi
    elif [ "$_dld" = wget ]; then
        if ! check_help_for wget --https-only --secure-protocol; then
            echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            wget "$1" -O "$2"
        else
            wget --https-only --secure-protocol=TLSv1_2 "$1" -O "$2"
        fi
    else
        err "Unknown downloader"   # should not reach here
    fi
}

check_help_for() {
    local _cmd
    local _arg
    local _ok
    _cmd="$1"
    _ok="y"
    shift

    # If we're running on OS-X, older than 10.13, then we always
    # fail to find these options to force fallback
    if check_cmd sw_vers; then
        if [ "$(sw_vers -productVersion | cut -d. -f2)" -lt 13 ]; then
            # Older than 10.13
            echo "Warning: Detected OS X platform older than 10.13"
           exit 1
        fi
    fi
}

main "$@" || exit 1