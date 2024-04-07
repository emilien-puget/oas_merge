package merge

import (
	"fmt"
	"reflect"

	"github.com/pb33f/libopenapi/datamodel/high"
	"github.com/pb33f/libopenapi/orderedmap"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func mergeComponents(mergedSpec *v3.Document, spec *v3.Document) (err error) {
	if spec.Components == nil {
		return nil
	}

	if mergedSpec.Components == nil {
		mergedSpec.Components = &v3.Components{}
	}

	// Merge schemas
	mergedSpec.Components.Schemas, err = mergeComponent(spec.Components.Schemas, mergedSpec.Components.Schemas)
	if err != nil {
		return fmt.Errorf("merge schemas: %w", err)
	}

	// Merge responses
	mergedSpec.Components.Responses, err = mergeComponent(spec.Components.Responses, mergedSpec.Components.Responses)
	if err != nil {
		return fmt.Errorf("merge responses: %w", err)
	}

	// Merge parameters
	mergedSpec.Components.Parameters, err = mergeComponent(spec.Components.Parameters, mergedSpec.Components.Parameters)
	if err != nil {
		return fmt.Errorf("merge parameters: %w", err)
	}

	// Merge examples
	mergedSpec.Components.Examples, err = mergeComponent(spec.Components.Examples, mergedSpec.Components.Examples)
	if err != nil {
		return fmt.Errorf("merge examples: %w", err)
	}

	// Merge request bodies
	mergedSpec.Components.RequestBodies, err = mergeComponent(spec.Components.RequestBodies, mergedSpec.Components.RequestBodies)
	if err != nil {
		return fmt.Errorf("merge request bodies: %w", err)
	}

	// Merge headers
	mergedSpec.Components.Headers, err = mergeComponent(spec.Components.Headers, mergedSpec.Components.Headers)
	if err != nil {
		return fmt.Errorf("merge headers: %w", err)
	}

	// Merge security schemes
	mergedSpec.Components.SecuritySchemes, err = mergeSecuritySchemes(spec.Components.SecuritySchemes, mergedSpec.Components.SecuritySchemes)
	if err != nil {
		return fmt.Errorf("merge security schemes: %w", err)
	}

	// Merge links
	mergedSpec.Components.Links, err = mergeComponent(spec.Components.Links, mergedSpec.Components.Links)
	if err != nil {
		return fmt.Errorf("merge links: %w", err)
	}

	// Merge callbacks
	mergedSpec.Components.Callbacks, err = mergeComponent(spec.Components.Callbacks, mergedSpec.Components.Callbacks)
	if err != nil {
		return fmt.Errorf("merge callbacks: %w", err)
	}

	return nil
}

func mergeSecuritySchemes(spec *orderedmap.Map[string, *v3.SecurityScheme], merged *orderedmap.Map[string, *v3.SecurityScheme]) (*orderedmap.Map[string, *v3.SecurityScheme], error) {
	if merged == nil {
		merged = orderedmap.New[string, *v3.SecurityScheme]()
	}
	for pair := spec.Oldest(); pair != nil; pair = pair.Next() {
		val, pres := merged.Get(pair.Key)
		if pres {
			if pair.Value.Type == "oauth2" {
				if pair.Value.Flows != nil {
					if pair.Value.Flows.Implicit != nil {
						mergeScopes(pair.Value.Flows.Implicit.Scopes, val.Flows.Implicit.Scopes)
					}
				}
			}
			if err := compareComponents(pair.Key, pair.Value, val); err != nil {
				return nil, err
			}
		} else {
			merged.Set(pair.Key, pair.Value)
		}
	}
	return merged, nil
}

func mergeScopes(spec, merge *orderedmap.Map[string, string]) {
	for pair := spec.Oldest(); pair != nil; pair = pair.Next() {
		merge.Set(pair.Key, pair.Value)
	}
}

func compareComponents(name string, sourceComponent, mergedComponent high.Renderable) error {
	if sourceComponent == nil {
		return nil
	}
	if mergedComponent == nil {
		// If the component does not exist in the merged spec, add it
		return nil
	}

	if reflect.DeepEqual(sourceComponent, mergedComponent) {
		return fmt.Errorf("component conflict: %s", name)
	}
	return nil
}

func mergeComponent[T high.Renderable](spec, merged *orderedmap.Map[string, T]) (*orderedmap.Map[string, T], error) {
	if merged == nil {
		merged = orderedmap.New[string, T]()
	}
	for pair := spec.Oldest(); pair != nil; pair = pair.Next() {
		val, pres := merged.Get(pair.Key)
		if pres {
			if err := compareComponents(pair.Key, pair.Value, val); err != nil {
				return nil, err
			}
		}
		merged.Set(pair.Key, pair.Value)
	}
	return merged, nil
}
