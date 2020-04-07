FROM golang:1.14-alpine

RUN apk add --no-cache gcc libc-dev gdb git && \
    git clone https://github.com/longld/peda.git ~/peda && \
    echo source ~/peda/peda.py >> ~/.gdbinit
