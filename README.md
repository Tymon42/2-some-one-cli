# Build on windows  
**Before building project, set the CGO_LDFLAGS and the CGO_CFLAGS environment variables, in order to make the Go build tools aware of the location of the VLC SDK files.**

```
set CGO_LDFLAGS=-L{your-sdk-path}
set CGO_CFLAGS=-I{your-sdk-path}\include
go build
```