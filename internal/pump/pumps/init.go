package pumps

var availablePumps map[string]Pump

func init() {
	availablePumps = make(map[string]Pump)

	availablePumps["csv"] = &CSVPump{}
	availablePumps["dummy"] = &DummyPump{}
}
