package html

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLi(t *testing.T) {
	tests := []struct {
		name string
		li   string
		want string
	}{
		{
			name: "li empty",
			want: "<li></li>",
		},
		{
			name: "li not empty",
			li:   "test",
			want: "<li>test</li>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlLi := Li(tt.li)

			assert.Equal(t, tt.want, htmlLi)
		})
	}
}

func TestUl(t *testing.T) {
	tests := []struct {
		name string
		ul   string
		want string
	}{
		{
			name: "ul empty",
			want: "<ul></ul>",
		},
		{
			name: "ul not empty",
			ul:   "test",
			want: "<ul>test</ul>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlUl := Ul(tt.ul)

			assert.Equal(t, tt.want, htmlUl)
		})
	}
}

func Test_document_AsString(t *testing.T) {
	type fields struct {
		start string
		body  string
		end   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty fields",
		},
		{
			name: "not empty fields",
			fields: fields{
				start: "1",
				body:  "2",
				end:   "3",
			},
			want: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &document{
				start: tt.fields.start,
				body:  tt.fields.body,
				end:   tt.fields.end,
			}

			assert.Equalf(t, tt.want, d.AsString(), "AsString()")
		})
	}
}
