module example

go 1.27

require (
	github.com/google/uuid v1.6.0
	github.com/insei/valigo v0.1.0
	github.com/insei/valigo/uuid v0.1.0
)

require (
	github.com/insei/fmap/v3 v3.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/insei/valigo => ../../
	github.com/insei/valigo/uuid => ../../uuid
)
