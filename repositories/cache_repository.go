package repositories

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"testapi/models"
	"testapi/utils"
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
	uniqueAddresses, _ := utils.RemoveDuplicateAddresses(allAddresses)

	// Write updated addresses back to the cache
	cacheContent, _ := json.Marshal(uniqueAddresses)
	file, err := os.Create(cacheFile)
	if err != nil {
		log.Printf("Error creating cache file %s: %v", cacheFile, err)
		return err
	}
	defer file.Close()

	_, err = file.Write(cacheContent)
	if err != nil {
		log.Printf("Error writing to cache file %s: %v", cacheFile, err)
		return err
	}

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
