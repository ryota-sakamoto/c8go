#!/bin/sh

echo "int one() { return 1; }" > one.c
echo "int two(int a, int b) { return a + b; }" > two.c

function run() {
    go run *.go "$@" > a.s
    if [ $? = 1 ]; then
        cat a.s
        return
    fi

    gcc -o a a.s one.c two.c
    ./a
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
        exit 1
    fi
}

function clean() {
    rm a.s a one.c two.c
}

check 0 "main() { 0; }"
check 42 "main() { 42; }"
check 21 "main() { 5+20-4; }"
check 41 "main() {  12 + 34 - 5 ; }"
check 47 "main() { 5+6*7; }"
check 15 "main() { 5*(9-6); }"
check 4 "main() { (3+5)/2; }"
check 10 "main() { -10+20; }"

check 1 "main() {  1  == 1; }"
check 0 "main() {  1  != 1; }"
check 0 "main() {  3  < 1; }"
check 0 "main() {  5  > 9; }"
check 1 "main() {  5  >= 4; }"
check 1 "main() {  4  <= 4; }"

check 3 "main() { \
a=3;a;\
}"
check 17 "main() { \
a = 3;\
b = a + 14;\
b;\
}"

check 20 "main() { foo = 1;\
bar = 2 + 17;\
foo + bar;\
}"
check 56 "main() { a=1;\
b = a + 27;\
return b * 2;\
}"

check 0 "main() { if (0) return 1;\
return 0;\
}"
check 6 "main() { a = 5; \
if (a == 5) a = a + 1;\
return a;\
}"

check 4 "main() { c = 2;\
if (c == 2) c = 4;\
else return 10;\
return c;\
}"

check 10 "main() { \
counter = 0;\
while (counter < 10) counter = counter + 1;\
return counter;\
}"

check 225 "main() { \
    a = 3;\
    b = 5;\
    if (1) {\
        a = a * b;\
        b = a * a;\
    }\
    return b;\
}"

check 4 "main() { return 3 + one(); }"
check 100 "main() { return two(1, 9) * two(6, 4); }"

check 9 "three() { return 3; } main() { a = three(); return a * three(); }"

clean