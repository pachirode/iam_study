package analytics

type AnalyticsFilters struct {
	Usernames        []string `json:"usernames"`
	SkippedUsernames []string `json:"skip_usernames"`
}

func (filters AnalyticsFilters) ShouldFilter(record AnalyticsRecord) bool {
	switch {
	case len(filters.SkippedUsernames) > 0 && stringInSlice(record.Username, filters.SkippedUsernames):
		return true
	case len(filters.Usernames) > 0 && !stringInSlice(record.Username, filters.Usernames):
		return true
	}

	return false
}

func (filters AnalyticsFilters) HasFilter() bool {
	if len(filters.SkippedUsernames) == 0 && len(filters.Usernames) == 0 {
		return false
	}

	return true
}
