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
	b.WriteString(`<TR><TD VALIGN="top" ROWSPAN="2">`)

	b.WriteString(`<A HREF="/` + string(p.ID) + `"><IMG WIDTH="100" SRC="` + p.Icon().URL + `"></A>`)
	b.WriteString(`<BR>`)
	b.WriteString(`<STRONG><a HREF="/` + string(p.ID) + `">` + p.Name() + `</a></STRONG>`)

	b.WriteString(`</TD><TD valign="top">`)
	b.Write(n.Content())

	attachments := n.Attachments()
	if len(attachments) > 0 {
		b.WriteString(`<TABLE BORDER="1" CELLPADDINg="5"><TBODY>`)
		b.WriteString(`<TR><TD COLSPAN="` + fmt.Sprintf("%d", len(attachments)) + `">Attachments</TD></TR><TR>`)
		for _, a := range attachments {
			b.WriteString(`<TD>`)
			if a.Type == "Document" {
				if strings.HasPrefix(a.MediaType, "image/") {
					b.WriteString(`<A TARGET="_blank" HREF="` + a.URL + `"><IMG WIDTH="150" SRC="` + a.URL + `"></A>`)
				} else if strings.HasPrefix(a.MediaType, "video/") {
					b.WriteString(`<VIDEO CONTROLS WIDTH="150">`)
					b.WriteString(`<SOURCE SRC="` + a.URL + `" TYPE="` + a.MediaType + `">`)
					b.WriteString(`</VIDEO>`)
				} else {
					b.WriteString(`<A TARGEt="_blank" HREF="` + a.URL + `">` + a.URL + `</A>`)
				}
			}
			b.WriteString(`</TD>`)
		}
		b.WriteString(`</TBODY></TABLE>`)
	}
	b.WriteString(`</TD></TR><TR><TD>`)
	b.WriteString(`<A HREF="/` + string(n.ID) + `">` + n.Published() + `</A>`)
	b.WriteString(`</TD></TR>`)
	b.Flush()
}
