# PDL
PDL (Parallel Downloader) is a CLI library that will download an URL.

## PDL-CLI
If you want to use this package as a cli
```go
go run cmd/pdl-cli/main.go
```

## Concept
PDL will split the files in several chunks and download it concurrently.  
It was an experiment to understant library like context, the concept of chan.

## Inspiration
Parrallel downloader inspired by a workshop in the berlin gophercon [original link](https://www.353solutions.com/c/pdl22/) by Miki Tebeka