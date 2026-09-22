package model

// Category classifies what KIND of anime-community event this is —
// separate from Tag (which fandom/genre it's about). Stored as plain
// text for the same reason as Role: adding a category later is a
// Go-side change only.
type Category string

const (
	CategoryConvention     Category = "convention"      // large multi-vendor/multi-panel event
	CategoryDoujinMarket   Category = "doujin_market"   // artist alley / fan-work sales
	CategoryScreening      Category = "screening"       // anime/movie screening
	CategoryCosplayContest Category = "cosplay_contest" // cosplay competition or gathering
	CategoryGameTournament Category = "game_tournament" // anime/visual-novel game tournament
	CategoryMeetup         Category = "meetup"          // casual fan meetup
)

// AllCategories lists every valid category, in display order. Used for
// input validation and for exposing the full set to the frontend (so it
// doesn't have to hardcode the list in two places).
func AllCategories() []Category {
	return []Category{
		CategoryConvention,
		CategoryDoujinMarket,
		CategoryScreening,
		CategoryCosplayContest,
		CategoryGameTournament,
		CategoryMeetup,
	}
}

// IsValidCategory reports whether c is one of the known categories.
func IsValidCategory(c Category) bool {
	for _, a := range AllCategories() {
		if c == a {
			return true
		}
	}
	return false
}
