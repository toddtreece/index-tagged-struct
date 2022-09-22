package index

import (
	"reflect"
	"strings"

	"github.com/blugelabs/bluge"
)

type ObjectSummary struct {
	Name        string `json:"name,omitempty" index:"name,keyword,sortable,storevalue,searchtermpositions"`
	Description string `json:"description,omitempty" index:"-"`
}

type DashboardObjectSummary struct {
	ObjectSummary `json:",inline" index:"extend"`
	Tags          []string `json:"tags,omitempty" index:"tags,keyword,aggregatable,storevalue,searchtermpositions"`
}

type IndexDocument struct {
	doc *bluge.Document
}

func NewIndexDocument(uid string) *IndexDocument {
	doc := bluge.NewDocument(uid)
	return &IndexDocument{doc: doc}
}

func (d *IndexDocument) Parse(object any) {
	t := reflect.TypeOf(object)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags := field.Tag.Get("index")

		if tags == "" || tags == "-" {
			continue
		}

		value := reflect.ValueOf(object).Field(i)

		if tags == "extend" {
			d.Parse(value.Interface())
			continue
		}

		d.addField(tags, value)
	}
}

func (d *IndexDocument) addField(t string, value reflect.Value) {
	tags := strings.Split(t, ",")
	name := tags[0]
	fieldType := tags[1]
	options := tags[2:]

	for _, field := range d.getFields(name, fieldType, options, value) {
		d.doc.AddField(field)
	}

}

func (d *IndexDocument) getFields(name, fieldType string, options []string, value reflect.Value) []*bluge.TermField {
	fields := make([]*bluge.TermField, 0)

	switch fieldType {
	case "keyword":
		fallthrough
	default:
		switch value.Kind() {
		case reflect.Slice:
			for _, v := range value.Interface().([]string) {
				fields = append(fields, bluge.NewKeywordField(name, v))
			}
		default:
			field := bluge.NewKeywordField(name, d.ensureValue(value).String())
			field = addFieldOptions(field, options)
			fields = append(fields, field)
		}
	}

	return fields
}

func addFieldOptions(field *bluge.TermField, options []string) *bluge.TermField {
	for _, option := range options {
		switch option {
		case "sortable":
			field = field.Sortable()
		case "storevalue":
			field = field.StoreValue()
		case "searchtermpositions":
			field = field.SearchTermPositions()
		case "aggregatable":
			field = field.Aggregatable()
		}
	}
	return field
}

func (d *IndexDocument) ToBlugeDocument() *bluge.Document {
	return d.doc
}

func (i *IndexDocument) ensureValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
