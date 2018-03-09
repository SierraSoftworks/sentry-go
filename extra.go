package sentry

import "encoding/json"

// Extra allows you to provide additional arbitrary metadata with your
// event. This data is not searchable, but can be invaluable in identifying
// the cause of a problem.
func Extra(extra map[string]string) Option {
	return &extraOption{extra}
}

type extraOption struct {
	extra map[string]string
}

func (o *extraOption) Class() string {
	return "extra"
}

func (o *extraOption) Merge(other Option) Option {
	if ot, ok := other.(*extraOption); ok {
		extra := make(map[string]string, len(o.extra))
		for k, v := range o.extra {
			extra[k] = v
		}

		for k, v := range ot.extra {
			extra[k] = v
		}

		return &extraOption{extra}
	}

	return other
}

func (o *extraOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.extra)
}
