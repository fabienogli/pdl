# PDL
PDL (Parallel Downloader) is a CLI library that will download an URL.

## Concept
PDL will split the files in several chunks and download it concurrently.  
It was an experiment to understand several standard libraries like context, http, the concept of chan and its error handling.


## PDL-CLI
If you want to use this package as a cli
```go
go run github.com/fabienogli/pdl/cmd/pdl-cli --help
```

## Inspiration
Parrallel downloader inspired by a workshop in the berlin gophercon [original link](https://www.353solutions.com/c/pdl22/) by Miki Tebeka