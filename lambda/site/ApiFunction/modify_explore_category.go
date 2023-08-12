package APIFunction

import (
	"encoding/json"
	"site/Data"
)


func (functions *APIFunctions) ModifyExploreCategory(explore_category_json string, is_add_or_remove_json string) (string, error) {
	var explore_category Data.ExploreCategory
	if err := json.Unmarshal([]byte(explore_category_json), &explore_category); err != nil {
		return "", err
	}
	var is_add_or_remove bool
	if err := json.Unmarshal([]byte(is_add_or_remove_json), &is_add_or_remove); err != nil {
		return "", err
	}
	if err := functions.DBRepo.ModifyExploreCategory(explore_category, is_add_or_remove); err != nil {
		return "", err
	}
	return "Success ModifyExploreCategory", nil
}
