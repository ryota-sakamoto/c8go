FROM golang:1.15.6-alpine

RUN apk add --no-cache gcc libc-dev gdb git make && \
    git clone https://github.com/longld/peda.git ~/peda && \
    echo source ~/peda/peda.py >> ~/.gdbinit
