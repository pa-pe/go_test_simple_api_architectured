package repositories

import (
	"encoding/json"
	"os"
	"path/filepath"

	"testapi/models"
)

type CacheRepository interface {
	UpdateCache(name, last string, newAddresses []models.Address) error
	LoadCache(name, last string) ([]models.Address, error)
}

type FileCacheRepository struct {
	CacheDir string
}

func NewFileCacheRepository(cacheDir string) *FileCacheRepository {
	return &FileCacheRepository{CacheDir: cacheDir}
}

func (r *FileCacheRepository) ensureCacheDir() {
	if _, err := os.Stat(r.CacheDir); os.IsNotExist(err) {
		os.Mkdir(r.CacheDir, 0755)
	}
}

func (r *FileCacheRepository) UpdateCache(name, last string, newAddresses []models.Address) error {
	r.ensureCacheDir()

	cacheKey := name + "_" + last
	cacheFile := filepath.Join(r.CacheDir, cacheKey+".json")

	// Load previous addresses from cache
	var cachedAddresses []models.Address
	if _, err := os.Stat(cacheFile); err == nil {
		fileContent, _ := os.Open(cacheFile)
		defer fileContent.Close()
		json.NewDecoder(fileContent).Decode(&cachedAddresses)
	}

	// Combine cached addresses with new addresses and remove duplicates
	allAddresses := append(cachedAddresses, newAddresses...)
	uniqueAddresses, _ := RemoveDuplicateAddresses(allAddresses)

	// Write updated addresses back to the cache
	cacheContent, _ := json.Marshal(uniqueAddresses)
	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(cacheContent)

	return nil
}

func (r *FileCacheRepository) LoadCache(name, last string) ([]models.Address, error) {
	cacheKey := name + "_" + last
	cacheFile := filepath.Join(r.CacheDir, cacheKey+".json")

	var cachedAddresses []models.Address
	if _, err := os.Stat(cacheFile); err == nil {
		fileContent, _ := os.Open(cacheFile)
		defer fileContent.Close()
		json.NewDecoder(fileContent).Decode(&cachedAddresses)
	}
	return cachedAddresses, nil
}

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
