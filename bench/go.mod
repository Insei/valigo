module bench

go 1.22

require (
	github.com/google/uuid v1.6.0
	github.com/insei/valigo v1.0.0
)

require github.com/insei/fmap/v3 v3.1.1 // indirect

replace github.com/insei/valigo v1.0.0 => ../
