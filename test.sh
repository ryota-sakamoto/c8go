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

check 0 "0;"
check 42 "42;"
check 21 "5+20-4;"
check 41 " 12 + 34 - 5 ;"
check 47 "5+6*7;"
check 15 "5*(9-6);"
check 4 "(3+5)/2;"
check 10 "-10+20;"

check 1 " 1  == 1;"
check 0 " 1  != 1;"
check 0 " 3  < 1;"
check 0 " 5  > 9;"
check 1 " 5  >= 4;"
check 1 " 4  <= 4;"

check 3 "a=3;a;"
check 17 "a = 3; b = a + 14; b;"

check 20 "foo = 1; bar = 2 + 17; foo + bar;"
check 56 "a=1; b = a + 27; return b * 2;"

clean