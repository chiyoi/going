usage() {
    echo "scripts:"
    echo "./scripts run"
    echo "./scripts install"
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