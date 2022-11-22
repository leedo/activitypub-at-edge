package render

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/leedo/activitypub-at-edge/activitypub"
)

func Note(w io.Writer, p *activitypub.Person, n *activitypub.Note) {
	b := bufio.NewWriter(w)
	b.WriteString(`<tr><td valign="top" rowspan="2">`)

	b.WriteString(`<a href="/` + string(p.ID) + `"><img width="100" src="` + p.Icon().URL + `"></a>`)
	b.WriteString(`<br>`)
	b.WriteString(`<strong><a href="/` + string(p.ID) + `">` + p.Name() + `</a></strong>`)

	b.WriteString(`</td><td valign="top">`)
	b.Write(n.Content())

	attachments := n.Attachments()
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
	b.WriteString(`<a href="/` + string(n.ID) + `">` + n.Published() + `</a>`)
	b.WriteString(`</td></tr>`)
	b.Flush()
}
