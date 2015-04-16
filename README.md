# speed

Speed measurement in golang

## install

```
go get github.com/go-torrent/speed
```

## usage

``` golang
import (
	"github.com/go-torrent/speed"

	"crypto/rand"
	"fmt"
)

func main() {
	b := make([]byte, 1024)
	g := speed.NewGauge()
	r, _ := rand.Read(b)

	speed := g.Progress(r)

	fmt.Printf("%v bytes/second", speed)
}
```

## License

MIT

This is a port from [speedometer](https://github.com/mafintosh/speedometer)
which is also MIT licensed.