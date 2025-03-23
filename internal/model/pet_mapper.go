package model

import "encoding/json"

func PetDBToPet(petDB PetDB) (Pet, error) {
	var pet Pet

	pet.ID = petDB.ID
	pet.Name = petDB.Name
	pet.Status = petDB.Status

	if err := json.Unmarshal([]byte(petDB.PhotoUrls), &pet.PhotoUrls); err != nil {
		return pet, err
	}
	if err := json.Unmarshal([]byte(petDB.Tags), &pet.Tags); err != nil {
		return pet, err
	}
	if err := json.Unmarshal([]byte(petDB.Category), &pet.Category); err != nil {
		return pet, err
	}

	return pet, nil
}

func PetToPetDB(pet Pet) (PetDB, error) {
	var petDB PetDB

	petDB.ID = pet.ID
	petDB.Name = pet.Name
	petDB.Status = pet.Status

	photoUrlsJSON, err := json.Marshal(pet.PhotoUrls)
	if err != nil {
		return petDB, err
	}
	tagsJSON, err := json.Marshal(pet.Tags)
	if err != nil {
		return petDB, err
	}
	categoryJSON, err := json.Marshal(pet.Category)
	if err != nil {
		return petDB, err
	}

	petDB.PhotoUrls = string(photoUrlsJSON)
	petDB.Tags = string(tagsJSON)
	petDB.Category = string(categoryJSON)

	return petDB, nil
}
