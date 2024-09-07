package utils

import (
	"testapi/models"
	"testing"
)

func TestRemoveDuplicateAddresses(t *testing.T) {
	addresses := []models.Address{
		{Country: "USA", City: "New York"},
		{Country: "USA", City: "New York"}, // Дубликат
		{Country: "Canada", City: "Toronto"},
	}

	expectedLength := 2
	uniqueAddresses, duplicatesRemoved := RemoveDuplicateAddresses(addresses)

	if len(uniqueAddresses) != expectedLength {
		t.Errorf("Expected %d unique addresses, got %d", expectedLength, len(uniqueAddresses))
	}

	if duplicatesRemoved != 1 {
		t.Errorf("Expected 1 duplicate removed, got %d", duplicatesRemoved)
	}
}
