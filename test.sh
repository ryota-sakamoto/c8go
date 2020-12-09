#!/bin/sh

mkdir -p tmp
CURRENT_DIR=$(pwd)

echo "int one() { return 1; }" > tmp/one.c
echo "int two(int a, int b) { return a + b; }" > tmp/two.c
cat <<EOF > tmp/p.c
#include <stdio.h>
void p(int v) { printf("%d\n", v); }
EOF
cat <<EOF > tmp/alloc4.c
#include <stdlib.h>
void alloc4(int **base, int a, int b, int c, int d) {
    int *pointer = (int *)malloc(sizeof(int)*4);
    *base = pointer;
    *(pointer) = a;
    *(pointer+1) = b;
    *(pointer+2) = c;
    *(pointer+3) = d;
}
EOF

function run() {
    cd "$CURRENT_DIR"/tmp

    ../bin/c8go "$@" > a.s
    if [ $? = 1 ]; then
        cat a.s
        return
    fi

    gcc -g -O0 -o a a.s one.c two.c p.c alloc4.c
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

check 0 "int main() { 0; }"
check 42 "int main() { 42; }"
check 21 "int main() { 5+20-4; }"
check 41 "int main() {  12 + 34 - 5 ; }"
check 47 "int main() { 5+6*7; }"
check 15 "int main() { 5*(9-6); }"
check 4 "int main() { (3+5)/2; }"
check 10 "int main() { -10+20; }"

check 1 "int main() {  1  == 1; }"
check 0 "int main() {  1  != 1; }"
check 0 "int main() {  3  < 1; }"
check 0 "int main() {  5  > 9; }"
check 1 "int main() {  5  >= 4; }"
check 1 "int main() {  4  <= 4; }"

check 3 "int main() { \
int a;\
a=3;a;\
}"
check 17 "int main() { \
int a;\
int b;\
a = 3;\
b = a + 14;\
b;\
}"

check 20 "int main() { int foo;\
int bar;\
foo = 1;\
bar = 2 + 17;\
foo + bar;\
}"
check 56 "int main() { \
int a;\
int b;\
a=1;\
b = a + 27;\
return b * 2;\
}"

check 0 "int main() { if (0) return 1;\
return 0;\
}"
check 6 "int main() { \
int a;\
a = 5; \
if (a == 5) a = a + 1;\
return a;\
}"

check 4 "int main() { \
int c;\
c = 2;\
if (c == 2) c = 4;\
else return 10;\
return c;\
}"

check 10 "int main() { \
int counter;\
counter = 0;\
while (counter < 10) counter = counter + 1;\
return counter;\
}"

check 225 "int main() { \
    int a;\
    int b;\
    a = 3;\
    b = 5;\
    if (1) {\
        a = a * b;\
        b = a * a;\
    }\
    return b;\
}"

check 4 "int main() { return 3 + one(); }"
check 100 "int main() { return two(1, 9) * two(6, 4); }"

check 9 "int three() { return 3; } int main() { int a; a = three(); return a * three(); }"

check 99 "int sum(int x, int y, int z) { return (x + y) * z; } int main() { return sum(10, 23, 3); }"

check 10 "int f(int x) { return x - 10; } int main() { x = 23; return f(x - 3); }"

check 120 "int f(int x) {\
    if (x == 1) return 1;\
    return f(x - 1) * x;\
}\
int main() {\
    f(5);\
}"

check 89 "int fib(int x) {\
    if (x == 0) return 1;\
    if (x == 1) return 1;\
    return fib(x - 1) + fib(x - 2);\
}\
int main() {\
    fib(10);\
}"

check 3 "int main() {\
    int x;\
    int y;\
    int z;\
    x = 3;\
    y = 5;\
    z = &y + 8;\
    return *z;\
}"

check 3 "int main() {\
    int x;\
    int *y;\
    y = &x;\
    *y = 3;\
    return *y;\
}"

check 8 "int main() {\
    int *a;\
    alloc4(&a, 1, 2, 4, 8);\
    int *b;\
    b = a + 4;\
    b = b - 1;\
    return *b;\
}"

check 4 "int main() {\
    int a;\
    return sizeof(a);\
}"
