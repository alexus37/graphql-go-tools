package postprocess

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

func TestCreateMultiNodes_ProcessFetchTree(t *testing.T) {
	t.Run("no parallel nodes", func(t *testing.T) {
		processor := &createMultiNodes{}
		input := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}}}),
		)
		processor.ProcessFetchTree(input)
		expected := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}}}),
		)
		require.Equal(t, expected, input)
	})

	t.Run("parallel nodes but no entity fetches", func(t *testing.T) {
		processor := &createMultiNodes{}
		ds := resolve.FetchInfo{
			DataSourceID:   "1",
			DataSourceName: "ds1",
		}
		input := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		processor.ProcessFetchTree(input)
		expected := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		require.Equal(t, expected, input)
	})

	t.Run("parallel nodes with single entity fetches", func(t *testing.T) {
		processor := &createMultiNodes{}
		ds := resolve.FetchInfo{
			DataSourceID:   "1",
			DataSourceName: "ds1",
		}
		input := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		processor.ProcessFetchTree(input)
		expected := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		require.Equal(t, expected, input)
	})

	t.Run("parallel nodes with two entity fetches but to different datasources", func(t *testing.T) {
		processor := &createMultiNodes{}
		ds1 := resolve.FetchInfo{
			DataSourceID:   "1",
			DataSourceName: "ds1",
		}
		ds2 := resolve.FetchInfo{
			DataSourceID:   "2",
			DataSourceName: "ds2",
		}
		input := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds2,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		processor.ProcessFetchTree(input)
		expected := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds2,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		require.Equal(t, expected, input)

	})

	t.Run("parallel nodes with two entity fetches but to the same datasources", func(t *testing.T) {
		processor := &createMultiNodes{}
		ds1 := resolve.FetchInfo{
			DataSourceID:   "1",
			DataSourceName: "ds1",
		}

		input := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		processor.ProcessFetchTree(input)
		expected := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Multi(
					[]*resolve.FetchTreeNode{
						resolve.Single(&resolve.EntityFetch{
							FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
							Info:              &ds1,
						}),
						resolve.Single(&resolve.EntityFetch{
							FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
							Info:              &ds1,
						}),
					},
				),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		require.Equal(t, expected, input)
	})

	t.Run("parallel nodes with two entity fetches but to the same datasources and a single fetch", func(t *testing.T) {
		processor := &createMultiNodes{}
		ds1 := resolve.FetchInfo{
			DataSourceID:   "1",
			DataSourceName: "ds1",
		}

		input := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
				resolve.Single(&resolve.EntityFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 7, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		processor.ProcessFetchTree(input)
		expected := resolve.Sequence(
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 0}}),
			resolve.Parallel(
				resolve.Single(&resolve.SingleFetch{
					FetchDependencies: resolve.FetchDependencies{FetchID: 7, DependsOnFetchIDs: []int{0}},
					Info:              &ds1,
				}),
				resolve.Multi(
					[]*resolve.FetchTreeNode{
						resolve.Single(&resolve.EntityFetch{
							FetchDependencies: resolve.FetchDependencies{FetchID: 1, DependsOnFetchIDs: []int{0}},
							Info:              &ds1,
						}),
						resolve.Single(&resolve.EntityFetch{
							FetchDependencies: resolve.FetchDependencies{FetchID: 2, DependsOnFetchIDs: []int{0}},
							Info:              &ds1,
						}),
					},
				),
			),
			resolve.Single(&resolve.SingleFetch{FetchDependencies: resolve.FetchDependencies{FetchID: 3, DependsOnFetchIDs: []int{1}}}),
		)
		require.Equal(t, expected, input)
	})
}
