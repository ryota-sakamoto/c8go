#!/bin/bash

function run() {
    go run *.go "$@" > a.s
    docker run -v $(pwd):/home -w /home --rm gcc-image gcc -o a a.s
    docker run -v $(pwd):/home -w /home --rm gcc-image /home/a
}

function check() {
    expected="$1"
    input="$2"

    run "$input"
    actual="$?"

    if [ "$expected" = "$actual" ]; then
        echo "$input => $actual"
    else
        echo "$input => $actual, but want $expected"
    fi
}

function clean() {
    rm a.s a
}

check 0 0
check 42 42
check 21 "5+20-4"
check 41 " 12 + 34 - 5 "

clean