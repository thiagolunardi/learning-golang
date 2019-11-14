package models

import (
	"github.com/matryer/is"
	"testing"
)

func TestSetAsDone(t *testing.T) {
	// arrange
	is := is.New(t)

	item := Item{
		ID:    1,
		Done:  false,
		Title: "Test",
	}

	// act
	beforeAct := item.Done
	item.SetAsDone()
	afterAct := item.Done

	// assert
	is.True(!beforeAct)
	is.True(afterAct)
	is.Equal(item.ID, 1)
	is.Equal(item.Title, "Test")
}
