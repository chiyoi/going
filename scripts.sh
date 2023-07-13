usage() {
    echo "scripts:"
    echo "scripts.sh run"
    echo "    Test run."
    echo "scripts.sh install"
    echo "    Install to the default go binary directory."
}

run() {
    go run . $@
}

install() {
    go install .
}

if test -z "$1" -o -n "$(echo "$1" | grep -Ex '\-{1,2}h(elp)?')"; then
usage
exit
fi

$@