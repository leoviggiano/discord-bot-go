package config

const (
	asth = "454894511852879877"
	bolo = "143928595755040768"
)

var adminList = []string{asth, bolo}

func Admin(id string) bool {
	for _, admin := range adminList {
		if id == admin {
			return true
		}
	}

	return false
}
