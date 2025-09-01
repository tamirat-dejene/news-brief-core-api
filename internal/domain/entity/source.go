package entity

// Source represents a news outlet in the database.
// It maps directly to a document in the 'sources' collection.

type LanguageType string

const (
	LanguageEN LanguageType = "en"
	LanguageAM LanguageType = "am"
)

type Source struct {
	ID               string       `bson:"_id,omitempty" json:"id"`
	Slug             string       `bson:"slug" json:"slug"`
	Name             string       `bson:"name" json:"name"`
	Description      string       `bson:"description" json:"description"`
	URL              string       `bson:"url" json:"url"`
	LogoURL          string       `bson:"logo_url" json:"logo_url"`
	Languages        LanguageType `bson:"languages" json:"languages"`
	Topics           []string     `bson:"topics" json:"topics"`
	ReliabilityScore float64      `bson:"reliability_score" json:"reliability_score"`
}

func SetLanguageType(lang string) LanguageType {
	switch lang {
	case "en":
		return LanguageEN
	case "am":
		return LanguageAM
	default:
		return LanguageEN
	}
}
