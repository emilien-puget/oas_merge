package merge

import (
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func Test_mergeOpenAPISpecs_simple(t *testing.T) {
	mainDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
`))
	require.NoError(t, err)
	mainModel, errs := mainDocument.BuildV3Model()
	require.Empty(t, errs)

	petsDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://petstore.swagger.io/v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response

`))
	require.NoError(t, err)
	petsModel, errs := petsDocument.BuildV3Model()
	require.Empty(t, errs)
	humansDocument, err := libopenapi.NewDocument([]byte(`
openapi: "3.1.0"
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://humanstore.swagger.io/v1
paths:
  /humans:
   get:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
`))
	require.NoError(t, err)
	humansModel, errs := humansDocument.BuildV3Model()
	require.Empty(t, errs)

	got, err := OpenAPISpecsWithMain(&mainModel.Model, []*v3.Document{&petsModel.Model, &humansModel.Model})
	assert.NoError(t, err)

	assert.Equal(t, `openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response
      servers:
        - url: http://petstore.swagger.io/v1
  /humans:
    get:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
      servers:
        - url: http://humanstore.swagger.io/v1
`, string(got.RenderWithIndention(2)))
}

func Test_mergeOpenAPISpecs_simple_components(t *testing.T) {
	mainDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
`))
	require.NoError(t, err)
	mainModel, errs := mainDocument.BuildV3Model()
	require.Empty(t, errs)

	petsDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://petstore.swagger.io/v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response

`))
	require.NoError(t, err)
	petsModel, errs := petsDocument.BuildV3Model()
	require.Empty(t, errs)
	humansDocument, err := libopenapi.NewDocument([]byte(`
openapi: "3.1.0"
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://humanstore.swagger.io/v1
paths:
  /humans:
   get:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
`))
	require.NoError(t, err)
	humansModel, errs := humansDocument.BuildV3Model()
	require.Empty(t, errs)

	got, err := OpenAPISpecsWithMain(&mainModel.Model, []*v3.Document{&petsModel.Model, &humansModel.Model})
	assert.NoError(t, err)

	assert.Equal(t, `openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response
      servers:
        - url: http://petstore.swagger.io/v1
  /humans:
    get:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
      servers:
        - url: http://humanstore.swagger.io/v1
`, string(got.RenderWithIndention(2)))
}

func Test_mergeOpenAPISpecs_common_path(t *testing.T) {
	mainDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
`))
	require.NoError(t, err)
	mainModel, errs := mainDocument.BuildV3Model()
	require.Empty(t, errs)

	petsDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://petstore.swagger.io/v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response

`))
	require.NoError(t, err)
	petsModel, errs := petsDocument.BuildV3Model()
	require.Empty(t, errs)
	humansDocument, err := libopenapi.NewDocument([]byte(`
openapi: "3.1.0"
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://humanstore.swagger.io/v1
paths:
  /pets:
   delete:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
`))
	require.NoError(t, err)
	humansModel, errs := humansDocument.BuildV3Model()
	require.Empty(t, errs)

	got, err := OpenAPISpecsWithMain(&mainModel.Model, []*v3.Document{&petsModel.Model, &humansModel.Model})
	assert.NoError(t, err)

	assert.Equal(t, `openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
paths:
  /pets:
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response
      servers:
        - url: http://petstore.swagger.io/v1
    delete:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
      servers:
        - url: http://humanstore.swagger.io/v1
`, string(got.RenderWithIndention(2)))
}

func Test_mergeOpenAPISpecs_path_parameters(t *testing.T) {
	mainDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
`))
	require.NoError(t, err)
	mainModel, errs := mainDocument.BuildV3Model()
	require.Empty(t, errs)

	petsDocument, err := libopenapi.NewDocument([]byte(`
openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://petstore.swagger.io/v1
paths:
  /pets:
    parameters:
      - name: titi
        in: query
        description: titi
        required: false
        schema:
          type: string
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response

`))
	require.NoError(t, err)
	petsModel, errs := petsDocument.BuildV3Model()
	require.Empty(t, errs)
	humansDocument, err := libopenapi.NewDocument([]byte(`
openapi: "3.1.0"
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://humanstore.swagger.io/v1
paths:
  /pets:
    parameters:
      - name: toto
        in: query
        description: toto
        required: false
        schema: 
          type: string
    delete:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
`))
	require.NoError(t, err)
	humansModel, errs := humansDocument.BuildV3Model()
	require.Empty(t, errs)

	got, err := OpenAPISpecsWithMain(&mainModel.Model, []*v3.Document{&petsModel.Model, &humansModel.Model})
	assert.NoError(t, err)

	assert.Equal(t, `openapi: 3.1.0
info:
  version: 1.0.0
  title: Swagger Petstore
  license:
    name: MIT
servers:
  - url: http://maintore.swagger.io/v1
paths:
  /pets:
    parameters:
      - name: titi
        in: query
        description: titi
        required: false
        schema:
          type: string
      - name: toto
        in: query
        description: toto
        required: false
        schema:
          type: string
    get:
      summary: List all pets
      operationId: listPets
      responses:
        default:
          description: default response
      servers:
        - url: http://petstore.swagger.io/v1
    delete:
      summary: List all humans
      operationId: listHumans
      responses:
        default:
          description: default response
      servers:
        - url: http://humanstore.swagger.io/v1
`, string(got.RenderWithIndention(2)))
}
