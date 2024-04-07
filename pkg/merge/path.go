package merge

import (
	"fmt"

	"github.com/pb33f/libopenapi/orderedmap"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func mergePaths(paths *v3.Paths, spec *v3.Document, mergedSpec *v3.Document) (*v3.Document, error) {
	for pair := paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
		itemClone := pair.Value

		for _, operation := range []*v3.Operation{
			itemClone.Get,
			itemClone.Put,
			itemClone.Post,
			itemClone.Delete,
			itemClone.Options,
			itemClone.Head,
			itemClone.Patch,
			itemClone.Trace,
		} {
			if operation != nil {
				operation.Servers = spec.Servers
			}
		}

		if mergedSpec.Paths == nil {
			mergedSpec.Paths = new(v3.Paths)
		}
		if mergedSpec.Paths.PathItems == nil {
			mergedSpec.Paths.PathItems = orderedmap.New[string, *v3.PathItem]()
		}
		if existingItem, ok := mergedSpec.Paths.PathItems.Get(pair.Key); ok {
			err := mergePathParameter(itemClone, existingItem)
			if err != nil {
				return nil, err
			}

			for _, operationName := range []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"} {
				var newOperation *v3.Operation
				switch operationName {
				case "get":
					newOperation = itemClone.Get
					existingOperation := existingItem.Get
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Get = newOperation
					}
				case "put":
					newOperation = itemClone.Put
					existingOperation := existingItem.Put
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Put = newOperation
					}
				case "post":
					newOperation = itemClone.Post
					existingOperation := existingItem.Post
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Post = newOperation
					}
				case "delete":
					newOperation = itemClone.Delete
					existingOperation := existingItem.Delete
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Delete = newOperation
					}
				case "options":
					newOperation = itemClone.Options
					existingOperation := existingItem.Options
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Options = newOperation
					}
				case "head":
					newOperation = itemClone.Head
					existingOperation := existingItem.Head
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Head = newOperation
					}
				case "patch":
					newOperation = itemClone.Patch
					existingOperation := existingItem.Patch
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Patch = newOperation
					}
				case "trace":
					newOperation = itemClone.Trace
					existingOperation := existingItem.Trace
					if existingOperation != nil && newOperation != nil {
						return nil, fmt.Errorf("Operation '%s' already exists.", operationName)
					} else if newOperation != nil {
						existingItem.Trace = newOperation
					}
				default:
					return nil, fmt.Errorf("Unknown operation '%s'", operationName)
				}
			}
		} else {
			mergedSpec.Paths.PathItems.Set(pair.Key, pair.Value)
		}
	}
	return nil, nil
}

func mergePathParameter(itemClone *v3.PathItem, existingItem *v3.PathItem) error {
	if len(itemClone.Parameters) > 0 && len(existingItem.Parameters) > 0 {
		itemClonePathparameters := orderedmap.New[string, *v3.Parameter]()
		for i := range itemClone.Parameters {
			itemClonePathparameters.Set(itemClone.Parameters[i].Name, itemClone.Parameters[i])
		}
		existingPathparameters := orderedmap.New[string, *v3.Parameter]()
		for i := range existingItem.Parameters {
			existingPathparameters.Set(existingItem.Parameters[i].Name, existingItem.Parameters[i])
		}
		newPathParameters, err := mergeComponent(itemClonePathparameters, existingPathparameters)
		if err != nil {
			return err
		}
		existingItem.Parameters = make([]*v3.Parameter, 0)
		for pathPair := newPathParameters.Oldest(); pathPair != nil; pathPair = pathPair.Next() {
			existingItem.Parameters = append(existingItem.Parameters, pathPair.Value)
		}
	}
	return nil
}
