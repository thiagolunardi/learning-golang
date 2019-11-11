package models

// Item is an object
type Item struct {
	// ID identification	
	ID    int			`json:"id"`
	// Done status of task
	Done  bool		`json:"done"`
	// Title of the task
	Title string	`json:"title"`
}

// Items is an object
type Items []Item

// DataSeed return sample data
func DataSeed() Items {
	return Items{
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