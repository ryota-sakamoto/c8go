#!/bin/sh

echo "int one() { return 1; }" > one.c
echo "int two(int a, int b) { return a + b; }" > two.c
cat <<EOF > p.c
#include <stdio.h>
void p(int v) { printf("%d\n", v); }
EOF

function run() {
    ./c8go "$@" > a.s
    if [ $? = 1 ]; then
        cat a.s
        return
    fi

    gcc -g -O0 -o a a.s one.c two.c p.c
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

check 99 "sum(x, y, z) { return (x + y) * z; } main() { return sum(10, 23, 3); }"

check 10 "f(x) { return x - 10; } main() { x = 23; return f(x - 3); }"

check 120 "f(x) {\
    if (x == 1) return 1;\
    return f(x - 1) * x;\
}\
main() {\
    f(5);\
}"

check 89 "fib(x) {\
    if (x == 0) return 1;\
    if (x == 1) return 1;\
    return fib(x - 1) + fib(x - 2);\
}\
main() {\
    fib(10);\
}"

check 3 "main() {\
    x = 3;\
    y = 5;\
    z = &y + 8;\
    return *z;\
}"
