package merge

import v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

// OpenAPISpecsWithMain merges the specs using main for the top level information.
func OpenAPISpecsWithMain(main *v3.Document, specs []*v3.Document) (*v3.Document, error) {
	mergedSpec := &v3.Document{
		Info:    main.Info,
		Version: main.Version,
		Servers: main.Servers,
	}

	for _, spec := range specs {
		paths := spec.Paths
		if paths != nil {
			document, err := mergePaths(paths, spec, mergedSpec)
			if err != nil {
				return document, err
			}
		}

		err := mergeComponents(mergedSpec, spec)
		if err != nil {
			return nil, err
		}
	}

	return mergedSpec, nil
}
