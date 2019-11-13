package models

// Item is an object
type Item struct {
	// ID identification
	ID int `json:"id"`
	// Done status of task
	Done bool `json:"done"`
	// Title of the task
	Title string `json:"title"`
}

// SetAsDone -
func (item *Item) SetAsDone() {
	item.Done = true
}

// Items is a collection
type Items []Item

// Main constructor
func Main() *Items {
	return &Items{}
}

// DataSeed return sample data
func DataSeed() []Item {
	return []Item{
		Item{
			ID:    1,
			Done:  false,
			Title: "Item A",
		},
		Item{
			ID:    2,
			Done:  false,
			Title: "Item B",
		},
	}
}
