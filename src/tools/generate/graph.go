// Copyright 2024 The Hexya Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

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
	return &ModelGraph{
		Nodes: make(map[string]*ModelNode),
	}
}

// AddNode adds a new node to the graph if it does not exist yet
func (g *ModelGraph) AddNode(name string) *ModelNode {
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
	model.Mixins = append(model.Mixins, mixin)
}

// AddEmbedEdge adds an embed edge from model to embed
func (g *ModelGraph) AddEmbedEdge(model *ModelNode, embed *ModelNode) {
	model.Embeds = append(model.Embeds, embed)
}

// Inflate populates each node of the graph with the fields and methods
// of its dependencies (mixins and embeds).
func (g *ModelGraph) Inflate() {
	for _, node := range g.Nodes {
		g.inflateNode(node)
	}
}

// inflateNode recursively inflates the given node with the data of its dependencies.
func (g *ModelGraph) inflateNode(node *ModelNode) {
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
