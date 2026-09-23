package uws1

// Components holds reusable objects scoped to the UWS document.
type Components struct {
	// Variables is an open-shape map of JSON-compatible values. Keys follow the
	// historical componentNamePattern before UWS 1.10; 1.10+ additionally
	// requires expression-addressable names.
	Variables  map[string]any `json:"variables,omitempty" yaml:"variables,omitempty" hcl:"variables,optional"`
	Extensions map[string]any `json:"-" yaml:"-" hcl:"extensions,block"`
}

type componentsAlias Components

var componentsKnownFields = []string{
	"variables",
}

func (c *Components) UnmarshalJSON(data []byte) error {
	var alias componentsAlias
	_, extensions, err := unmarshalCoreWithExtensions(data, "components", componentsKnownFields, &alias)
	if err != nil {
		return err
	}
	*c = Components(alias)
	c.Extensions = extensions
	return nil
}

func (c Components) MarshalJSON() ([]byte, error) {
	alias := componentsAlias(c)
	return marshalWithExtensions(&alias, c.Extensions)
}
