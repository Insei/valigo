module github.com/insei/valigo/integrations/goi18n

go 1.27

replace github.com/insei/valigo => ../../

require (
	github.com/insei/valigo v0.2.0
	github.com/nicksnyder/go-i18n/v2 v2.6.1
	github.com/stretchr/testify v1.12.1
	golang.org/x/text v0.43.0
)

require (
	github.com/insei/fmap/v3 v3.3.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
