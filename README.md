rfw is a very simple tool to watch for changes in one of more paths and run
a command whevener a change is detected.

This is mostly done for my own use.

Example:

```
rfw -path . -path ./something go run .
```

If anything changes in either `.` or `./something`, then `rfw` will execute `go
run .`

Note that `rfw` explicitly ignores `chmod`
