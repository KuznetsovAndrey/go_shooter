# ab-like tool, written in go

## My first expirience with golang. 

Usage:
go run main.go (or compile and execute)
params:
  - `-H string`
    	headers, separated by semicolon
  - `-X string`
    	method (default "GET")
  - `-c int`
    	parallel shoots in time (default 1)
  - `-d int`
    	delay in ms
  - `-h string`
    	host
  - `-n int`
    	count of shots (default 1)
  - `-p int`
    	port (default -1)
  - `-ssl int`
      use https schema (1 to use, 0 or skip for http)
