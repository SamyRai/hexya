package parser

import "fmt"

// A ModelGraph holds the graph of all models
type ModelGraph struct {
	Nodes map[string]*ModelNode
}

// A ModelNode is a node in the ModelGraph
type ModelNode struct {
	Name         string
	ModelType    string
	IsModelMixin bool
	Fields       map[string]FieldASTData
	Methods      map[string]MethodASTData
	Mixins       []*ModelNode
	Embeds       []*ModelNode
}

// NewModelGraph returns a new ModelGraph instance
func NewModelGraph() *ModelGraph {
	fmt.Println("NewModelGraph called")
	defer fmt.Println("NewModelGraph returned")
	return &ModelGraph{
		Nodes: make(map[string]*ModelNode),
	}
}

// AddNode adds a new node to the graph if it does not exist yet
func (g *ModelGraph) AddNode(name string) *ModelNode {
	fmt.Printf("AddNode called with name: %s\n", name)
	defer fmt.Println("AddNode returned")
	if node, ok := g.Nodes[name]; ok {
		return node
	}
	node := &ModelNode{
		Name:    name,
		Fields:  make(map[string]FieldASTData),
		Methods: make(map[string]MethodASTData),
	}
	g.Nodes[name] = node
	return node
}

// AddMixinEdge adds a mixin edge from model to mixin
func (g *ModelGraph) AddMixinEdge(model *ModelNode, mixin *ModelNode) {
	fmt.Printf("AddMixinEdge called with model: %s, mixin: %s\n", model.Name, mixin.Name)
	defer fmt.Println("AddMixinEdge returned")
	model.Mixins = append(model.Mixins, mixin)
}

// AddEmbedEdge adds an embed edge from model to embed
func (g *ModelGraph) AddEmbedEdge(model *ModelNode, embed *ModelNode) {
	fmt.Printf("AddEmbedEdge called with model: %s, embed: %s\n", model.Name, embed.Name)
	defer fmt.Println("AddEmbedEdge returned")
	model.Embeds = append(model.Embeds, embed)
}

// Inflate populates each node of the graph with the fields and methods
// of its dependencies (mixins and embeds).
func (g *ModelGraph) Inflate() {
	fmt.Println("Inflate called")
	defer fmt.Println("Inflate returned")
	for _, node := range g.Nodes {
		g.inflateNode(node)
	}
}

// inflateNode recursively inflates the given node with the data of its dependencies.
func (g *ModelGraph) inflateNode(node *ModelNode) {
	fmt.Printf("inflateNode called with node: %s\n", node.Name)
	defer fmt.Println("inflateNode returned")
	for _, mixin := range node.Mixins {
		g.inflateNode(mixin)
		for fieldName, field := range mixin.Fields {
			if fieldName == "ID" {
				continue
			}
			field.MixinField = true
			if _, exists := node.Fields[fieldName]; !exists {
				node.Fields[fieldName] = field
			}
		}
		for methodName, method := range mixin.Methods {
			method.ToDeclare = true
			if _, exists := node.Methods[methodName]; !exists {
				node.Methods[methodName] = method
			}
		}
	}
	for _, embed := range node.Embeds {
		g.inflateNode(embed)
		for fieldName, field := range embed.Fields {
			if _, exists := node.Fields[fieldName]; exists {
				continue
			}
			embeddedField := field
			embeddedField.EmbedField = true
			node.Fields[fieldName] = embeddedField
		}
	}
}
