package block

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Block interface {
	Attributes() map[string]Attribute
	AttachEvalContext(ctx *hcl.EvalContext)
	HasModuleBlock() bool
	GetModuleBlock() (Block, error)
	Type() string
	Labels() []string
	Range() Range
	GetFirstMatchingBlock(names ...string) Block
	GetBlock(name string) Block
	AllBlocks() Blocks
	GetBlocks(name string) Blocks
	GetAttributes() []Attribute
	GetAttribute(name string) Attribute
	GetNestedAttribute(name string) Attribute
	Reference() *Reference
	LocalName() string
	FullName() string
	UniqueName() string
	TypeLabel() string
	NameLabel() string
	Clone(index cty.Value) Block
	IsCountExpanded() bool
	HasChild(childElement string) bool
	MissingChild(childElement string) bool
	MissingNestedChild(childElement string) bool
	InModule() bool
	Label() string
	HasBlock(childElement string) bool
	IsResourceType(resourceType string) bool
	IsEmpty() bool
	Values() cty.Value
	Context() *hcl.EvalContext
	ReadLines() (lines []string, comments []string, err error)
	IsNil() bool
	IsNotNil() bool
}
