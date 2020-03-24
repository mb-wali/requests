package util

// CopyRequestDetails makes a one-level-deep copy of a map. For copying request details, we only need to go one level
// deep because this service doesn't need to modify anything below the top level of the map.
func CopyRequestDetails(requestDetails map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range requestDetails {
		copy[k] = v
	}
	return copy
}
