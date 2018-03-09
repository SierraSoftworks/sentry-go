package sentry

import "encoding/json"

// Modules allows you to specify the versions of various modules
// used by your application.
func Modules(moduleVersions map[string]string) Option {
	return &modulesOption{moduleVersions}
}

type modulesOption struct {
	moduleVersions map[string]string
}

func (o *modulesOption) Class() string {
	return "modules"
}

func (o *modulesOption) Merge(other Option) Option {
	if ot, ok := other.(*modulesOption); ok {
		moduleVersions := make(map[string]string, len(o.moduleVersions))
		for k, v := range o.moduleVersions {
			moduleVersions[k] = v
		}

		for k, v := range ot.moduleVersions {
			moduleVersions[k] = v
		}

		return &modulesOption{moduleVersions}
	}

	return other
}

func (o *modulesOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.moduleVersions)
}
