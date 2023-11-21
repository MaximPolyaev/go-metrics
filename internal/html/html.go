// Package html for create html5 document
package html

const (
	start = `
<!doctype html>
<html lang="ru">
<body>
`
	end = `
</body>
</html>
`
)

type Document struct {
	start string
	body  string
	end   string
}

func NewDocument() *Document {
	return &Document{
		start: start,
		end:   end,
	}
}

// SetBody - set body to document
func (d *Document) SetBody(body string) {
	d.body = body
}

// AsString - convert document to string
func (d *Document) AsString() string {
	return d.start + d.body + d.end
}

// Ul - ul tag for string
func Ul(ul string) string {
	return "<ul>" + ul + "</ul>"
}

// Li - li tag for string
func Li(li string) string {
	return "<li>" + li + "</li>"
}
