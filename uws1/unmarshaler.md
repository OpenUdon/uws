# UWS HCL Unmarshaler Ownership

`UnmarshalHCL` implementations follow the same ownership model as
`encoding/json.Unmarshaler`: the type that owns a field is responsible for
normalizing that field after the decoder has populated the local value.

The ownership rules are:

- A child struct should implement `UnmarshalHCL` if decoding that child directly with `dethcl.Unmarshal(childHCL, &child, labels...)` would fail or produce the wrong value.
- `Document` should implement `UnmarshalHCL` only if `Document` itself has fields Horizon cannot decode correctly by default.
- `Document` should not implement `UnmarshalHCL` as a workaround for broken children. Fix the child.
- Once every child is self-decoding where needed, parent decoding composes naturally.

Horizon invokes a type's custom HCL unmarshaller before falling back to default
field decoding. That means a parent does not need to walk the decoded object
tree to repair children that already implement their own local HCL behavior.
It should decode through an alias, assign the decoded value, and normalize only
fields owned by that type.

For `Document`, the current custom behavior is root-scoped: decode through
`documentHCLAlias`, normalize `Document.Variables`, and normalize root
`Document.Extensions`. Child structs own their own custom behavior for fields
such as dynamic maps, extension payloads, and escaped description text.
