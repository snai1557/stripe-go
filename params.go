package stripe

import (
	"fmt"
	"net/url"
	"strconv"
)

// ListParams is the structure that contains the common properties
// of any *ListParams structure.
type ListParams struct {
	Start, End string
	Limit      uint64
	Filters    Filters
}

// Params is the structure that contains the common properties
// of any *Params structure.
type Params struct {
	Exp         []string
	Meta        map[string]string
	AccessToken string
}

// Filters is a structure that contains a collection of filters for list-related APIs.
type Filters struct {
	f []*filter
}

// filter is the structure that contains a filter for list-related APIs.
// It ends up passing query string parameters in the format key[op]=value.
type filter struct {
	Key, Op, Val string
}

// Expand appends a new field to expand.
func (p *Params) Expand(f string) {
	p.Exp = append(p.Exp, f)
}

// AddFilter adds a new filter with a given key, op and value.
func (f *Filters) AddFilter(key, op, value string) {
	filter := &filter{Key: key, Op: op, Val: value}
	f.f = append(f.f, filter)
}

// appendTo adds the common parameters to the query string values.
func (p *Params) appendTo(body *url.Values) {
	for k, v := range p.Meta {
		body.Add(fmt.Sprintf("metadata[%v]", k), v)
	}

	for _, v := range p.Exp {
		body.Add("expand[]", v)
	}
}

// appendTo adds the common parameters to the query string values.
func (p *ListParams) appendTo(body *url.Values) {
	if len(p.Filters.f) > 0 {
		p.Filters.appendTo(body)
	}

	if len(p.Start) > 0 {
		body.Add("starting_after", p.Start)
	}

	if len(p.End) > 0 {
		body.Add("ending_before", p.End)
	}

	if p.Limit > 0 {
		if p.Limit > 100 {
			p.Limit = 100
		}

		body.Add("limit", strconv.FormatUint(p.Limit, 10))
	}
}

// appendTo adds the list of filters to the query string values.
func (f *Filters) appendTo(values *url.Values) {
	for _, v := range f.f {
		if len(v.Op) > 0 {
			values.Add(fmt.Sprintf("%v[%v]", v.Key, v.Op), v.Val)
		} else {
			values.Add(v.Key, v.Val)
		}
	}
}
