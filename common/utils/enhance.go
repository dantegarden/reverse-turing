package utils

func RemoveStringCreateNew(slice []string, str string) []string {
	var newSlice []string
	for _, v := range slice {
		if v != str {
			// 如果当前元素不是要删除的字符串，则添加到新切片
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
