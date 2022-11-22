package render

import (
	"bufio"
	"io"
	"strings"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Post(w io.Writer, p *activitypub.Person, o *activitypub.Object) {
	b := bufio.NewWriter(w)
	b.WriteString(`<tr><td>`)
	if p != nil {
		b.WriteString(`<img width="100" src="` + p.Icon().URL + `">`)
		b.WriteString(`<br>`)
		b.WriteString(`<strong>` + p.Name() + `</strong>`)
	}
	b.WriteString(`</td><td>`)
	if o != nil {
		b.Write(o.Content())
	} else {
		b.Write(o.ID)
	}
	for _, a := range o.Attachments() {
		if a.Type == "Document" {
			if strings.HasPrefix(a.MediaType, "image/") {
				b.WriteString(`<br>`)
				b.WriteString(`<img height="200" src="` + a.URL + `">`)
			} else if strings.HasPrefix(a.MediaType, "video/") {
				b.WriteString(`<video controls height="200">`)
				b.WriteString(`<source src="` + a.URL + `" type="` + a.MediaType + `">`)
				b.WriteString(`</video>`)
			}
		}
	}
	b.WriteString(`</td></tr>`)
	b.Flush()
}
