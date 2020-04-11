c8go
===

https://www.sigbus.info/compilerbook

## test

```
$ docker build . -t gcc-image
$ docker run -v $(pwd):/home -w /home --rm gcc-image make test
```
