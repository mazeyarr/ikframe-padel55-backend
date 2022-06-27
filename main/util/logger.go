package util

func GetLogPrefix(className, functionName string) string {
	return className + ":" + functionName
}
