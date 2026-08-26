package models

// DashboardStatsResponse is the response for dashboard statistics
// organizationCount is only visible to Site Admins
type DashboardStatsResponse struct {
	OrganizationCount *int64 `json:"organizationCount,omitempty"`
	GroupCount        int64  `json:"groupCount"`
}
