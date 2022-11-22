package render

import (
	"bufio"
	"io"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Person(w io.Writer, p *activitypub.Person) {
	b := bufio.NewWriter(w)
	b.WriteString(`<a href="/` + string(p.ID) + `"><img src="` + p.Image().URL + `" height="200"></a>`)
	b.WriteString(`<h1><a href="/` + string(p.ID) + `">` + p.Name() + `</a>'s Page</h1>`)
	b.Flush()
}
