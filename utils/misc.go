package utils

func Find(list []string, needle string) int {
	l := len(list)
	for i := 0; i < l; i++ {
		if list[i] == needle {
			return i
		}
	}
	return -1
}

func StringValue(data map[string]interface{}, key string, defaultValue string) string {
	if data[key] == nil {
		return defaultValue
	}
	value, ok := data[key].(string)
	if !ok {
		return defaultValue
	}
	if value == "" {
		return defaultValue
	}
	return value
}
