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

func (d *Document) SetBody(body string) {
	d.body = body
}

func (d *Document) AsString() string {
	return d.start + d.body + d.end
}

func Ul(ul string) string {
	return "<ul>" + ul + "</ul>"
}

func Li(li string) string {
	return "<li>" + li + "</li>"
}
