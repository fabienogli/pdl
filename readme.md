# PDL (Parallel Downloader)  

PDL is a CLI library for downloading files from a URL with parallel execution.  

## Concept  

PDL splits files into multiple chunks and downloads them concurrently. This project serves as an experiment to explore various Go standard libraries, including:  

- `context` for managing execution lifecycles  
- `http` for handling requests  
- Channels (`chan`) for concurrency and error handling  

## Using PDL-CLI  

You can use this package as a CLI tool:  

```sh
go run github.com/fabienogli/pdl/cmd/pdl-cli --help
```

## Inspiration
Parrallel downloader inspired by a workshop in the berlin gophercon [original link](https://www.353solutions.com/c/pdl22/) by Miki Tebeka