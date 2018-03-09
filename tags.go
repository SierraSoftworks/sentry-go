package sentry

import "encoding/json"

// Tags allow you to add additional tagging information to events which
// makes it possible to easily group and query similar events.
func Tags(tags map[string]string) Option {
	return &tagsOption{tags}
}

type tagsOption struct {
	tags map[string]string
}

func (o *tagsOption) Class() string {
	return "tags"
}

func (o *tagsOption) Merge(other Option) Option {
	if ot, ok := other.(*tagsOption); ok {
		tags := make(map[string]string, len(o.tags))
		for k, v := range o.tags {
			tags[k] = v
		}

		for k, v := range ot.tags {
			tags[k] = v
		}

		return &tagsOption{tags}
	}

	return other
}

func (o *tagsOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.tags)
}
