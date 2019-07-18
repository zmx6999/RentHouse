package utils

func StringValueEmpty(data map[string]interface{}, key string) bool {
	return data[key] == nil || data[key].(string) == ""
}

func GetStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if StringValueEmpty(data, key) {
		return defaultValue
	} else {
		return data[key].(string)
	}
}

func Find(list []string, str string) int {
	for i := 0; i < len(list); i++ {
		if list[i] == str {
			return i
		}
	}

	return -1
}
