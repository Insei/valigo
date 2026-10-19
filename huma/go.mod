module github.com/insei/valigo/huma

go 1.27

replace github.com/insei/valigo => ../

require (
	github.com/danielgtaylor/huma/v2 v2.39.1
	github.com/insei/fmap/v3 v3.3.0
	github.com/insei/valigo v0.1.0
	github.com/stretchr/testify v1.12.1
)

require (
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
