package utils

import "testapi/models"

func RemoveDuplicateAddresses(addresses []models.Address) ([]models.Address, int) {
	addressMap := make(map[string]models.Address)
	duplicatesRemoved := 0
	for _, addr := range addresses {
		key := addr.Country + "_" + addr.City
		if _, exists := addressMap[key]; exists {
			duplicatesRemoved++
		}
		addressMap[key] = addr
	}

	uniqueAddresses := []models.Address{}
	for _, addr := range addressMap {
		uniqueAddresses = append(uniqueAddresses, addr)
	}

	return uniqueAddresses, duplicatesRemoved
}
