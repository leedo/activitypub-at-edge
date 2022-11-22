package render

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Post(w io.Writer, p *activitypub.Person, o *activitypub.Object) {
	b := bufio.NewWriter(w)
	b.WriteString(`<tr><td valign="top" rowspan="2">`)
	if p != nil {
		b.WriteString(`<img width="100" src="` + p.Icon().URL + `">`)
		b.WriteString(`<br>`)
		b.WriteString(`<strong>` + p.Name() + `</strong>`)
	}
	b.WriteString(`</td><td valign="top">`)
	if o != nil {
		b.Write(o.Content())
	} else {
		b.Write(o.ID)
	}
	attachments := o.Attachments()
	if len(attachments) > 0 {
		b.WriteString(`<table border="1" cellpadding="5"><tbody>`)
		b.WriteString(`<tr><td colspan="` + fmt.Sprintf("%d", len(attachments)) + `">Attachments</td></tr><tr>`)
		for _, a := range attachments {
			b.WriteString(`<td>`)
			if a.Type == "Document" {
				if strings.HasPrefix(a.MediaType, "image/") {
					b.WriteString(`<a target="_blank" href="` + a.URL + `"><img width="150" src="` + a.URL + `"></a>`)
				} else if strings.HasPrefix(a.MediaType, "video/") {
					b.WriteString(`<video controls width="150">`)
					b.WriteString(`<source src="` + a.URL + `" type="` + a.MediaType + `">`)
					b.WriteString(`</video>`)
				} else {
					b.WriteString(`<a target="_blank" href="` + a.URL + `">` + a.URL + `</a>`)
				}
			}
			b.WriteString(`</td>`)
		}
		b.WriteString(`</tbody></table>`)
	}
	b.WriteString(`</td></tr><tr><td>`)
	b.WriteString(`<a href="/` + string(o.ID) + `">` + o.Published() + `</a>`)
	b.WriteString(`</td></tr>`)
	b.Flush()
}
