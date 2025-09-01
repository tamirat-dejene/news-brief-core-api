package entity

// BilingualField represents a field with both English and Amharic translations.
type BilingualField struct {
	EN string `bson:"en" json:"en"`
	AM string `bson:"am" json:"am"`
}

// Topic represents a news category.
type Topic struct {
	ID         string         `bson:"_id,omitempty" json:"id"`
	Slug       string         `bson:"slug" json:"slug"`
	Label      BilingualField `bson:"label" json:"label"`
	StoryCount int            `bson:"story_count" json:"story_count"`
}
