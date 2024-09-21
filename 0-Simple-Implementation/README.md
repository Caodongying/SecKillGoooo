# Initialize module
Execute:
1. ```go mod init 0-simple-implementation```
2. ```go mod tidy```
3. ``` go get github.com/gomodule/redigo/redis ```

好像每一次go get后都要run一下go mod tidy？🧐
如果一开始不go mod init的话，就无法执行go get，写的package main也不会被识别到。