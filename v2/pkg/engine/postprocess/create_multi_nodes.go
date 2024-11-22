package postprocess

import (
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

type createMultiNodes struct {
	disable bool
}

func (c *createMultiNodes) ProcessFetchTree(root *resolve.FetchTreeNode) {
	if c.disable {
		return
	}
	// iterate over all and look for parallel nodes
	for i := 0; i < len(root.ChildNodes); i++ {
		if root.ChildNodes[i].Kind == resolve.FetchTreeNodeKindParallel {
			c.handleParallelNode(root.ChildNodes[i])
		}
	}
}

func (c *createMultiNodes) handleParallelNode(node *resolve.FetchTreeNode) {
	// split the child nodes into groups that can be batched together by data source
	groups := make(map[string][]*resolve.FetchTreeNode)
	for _, child := range node.ChildNodes {
		dataSourceID := child.Item.Fetch.DataSourceInfo().ID
		if groups[dataSourceID] == nil {
			groups[dataSourceID] = make([]*resolve.FetchTreeNode, 0, 2)
		}
		groups[dataSourceID] = append(groups[dataSourceID], child)
	}

	var newChildNodes []*resolve.FetchTreeNode
	// iterate over the groups and create a multi nodes for each group wiht more then 2 nodes
	for _, group := range groups {
		if len(group) > 1 {
			// create a new multi node
			multiNode := resolve.Multi(group)
			newChildNodes = append(newChildNodes, multiNode)
		} else {
			newChildNodes = append(newChildNodes, group[0])
		}
	}

	// set the new child nodes for the parallel node
	node.ChildNodes = newChildNodes
}
