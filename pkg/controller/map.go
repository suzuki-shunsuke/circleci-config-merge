package controller

func mergeMap(base, child map[string]interface{}) map[string]interface{} {
	if base == nil {
		return child
	}
	for k, v := range child {
		base[k] = v
	}
	return base
}
