https://www.353solutions.com/c/pdl22/

why go ?
the free lunch is over --> processor doesn't go faster so the solution is hyperthreading
http://www.gotw.ca/publications/concurrency-ddj.htm

goroutines
c10k problem
https://www.packt.com/c10k-non-blocking-web-server-go/


context is used for timeout

go build --> dynamically linked binary
CGO_ENABLED=0 go build --> statically linked binary


go mod vendor --> allow to fix external libraries or remove dependencies from internet but it can be a pretty big folder

Fuzzing is a good practise cause we don't know very input
go test -fuzz=Fuzz

//effective go --> goto doc

only improve the performance after implementation, if there is a NEED (is 1sec enough)
for example make([]interface, nbChunks) > append(cs, chunk)
--> every time we append, we allocate a new slice, so there is more and more time to allocate 

go test -v only test the current package

Go doesn't wait for goroutine
-> if the main finished, it doesn't car
```go
for i := 0; i < 3; i++ {
    go func(){
        fmt.Println(i) // bug it only use the same i from the "i"
    }()
}
```
the first solution is to pass the value as a parameter
```go
for i := 0; i < 3; i++ {
    go func(val int){
        fmt.Println(val)
    }(i)
}
```
solution 2 is to create a loop var
```go
for i := 0; i < 3; i++ {
    i := i //Shadow i from the for loop
    go func(){
        fmt.Println(i)
    }()
}
```

Why ? -->
```go
// '{' create a new scope}
```

chan operations are blocking operations

fatalF should only be called in main

create a cli using the flag package

To list every format for go
go tool dist list
GOOS=dist go build ...


## Documentation
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/HEAD
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
