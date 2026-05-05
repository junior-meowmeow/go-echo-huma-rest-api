package database

func get[T any](m map[string]any, key string) T {
	val, _ := m[key].(T)
	return val
}

func getInt(m map[string]any, key string) int {
	if val, ok := m[key].(int32); ok {
		return int(val)
	}
	if val, ok := m[key].(int64); ok {
		return int(val)
	}
	return 0
}

func getStringOptional(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if valStr, ok2 := val.(string); ok2 {
			return valStr
		}
	}
	return ""
}
