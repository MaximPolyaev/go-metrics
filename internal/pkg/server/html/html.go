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

type Document interface {
	AsString() string
	SetBody(string)
}

type document struct {
	start string
	body  string
	end   string
}

func NewDocument() Document {
	return &document{
		start: start,
		end:   end,
	}
}

func (d *document) SetBody(body string) {
	d.body = body
}

func (d *document) AsString() string {
	return d.start + d.body + d.end
}

func Ul(ul string) string {
	return "<ul>" + ul + "</ul>"
}

func Li(li string) string {
	return "<li>" + li + "</li>"
}
