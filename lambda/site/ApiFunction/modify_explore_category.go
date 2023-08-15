package APIFunction

import (
	"encoding/json"
	"site/Data"
)

func (functions *APIFunctions) ModifyExploreCategory(explore_category_json string, modify_type string) (string, error) {
	var explore_category Data.ExploreCategory
	if err := json.Unmarshal([]byte(explore_category_json), &explore_category); err != nil {
		return "Error ModifyExploreCategory", err
	}
	if err := functions.DBRepo.ModifyExploreCategory(modify_type, explore_category); err != nil {
		return "Error ModifyExploreCategory", err
	}
	return "Success ModifyExploreCategory", nil
}
