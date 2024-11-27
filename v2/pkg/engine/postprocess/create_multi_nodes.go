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
		// select all fetches that are entity or batch entity fetches
		entityFetches := make([]*resolve.FetchTreeNode, 0, len(group))
		for _, child := range group {
			// check if the fetch is an entity or batch entity fetch
			switch child.Item.Fetch.(type) {
			case *resolve.EntityFetch, *resolve.BatchEntityFetch:
				entityFetches = append(entityFetches, child)
			default:
				newChildNodes = append(newChildNodes, child)
			}

		}

		if len(entityFetches) > 1 {
			// create a new multi node
			multiNode := resolve.Multi(entityFetches)
			newChildNodes = append(newChildNodes, multiNode)
		} else if len(entityFetches) == 1 {
			newChildNodes = append(newChildNodes, entityFetches[0])
		}
	}

	// set the new child nodes for the parallel node
	node.ChildNodes = newChildNodes
}
