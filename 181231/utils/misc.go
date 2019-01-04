package utils

func AddDomain2Url(url string) string {
	return "http://"+FastDFSHost+":"+FastDFSPort+"/"+url
}
