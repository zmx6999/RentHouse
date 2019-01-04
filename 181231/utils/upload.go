package utils

import (
	"net/http"
	"mime/multipart"
	"path"
	"errors"
	"strings"
)

func find(x string,a []string) int {
	for k,v:=range a{
		if v==x {
			return k
		}
	}
	return -1
}

func PrepareUpload(r *http.Request,key string,allowedTypes []string,allowedMaxSize int) ([]byte, *multipart.FileHeader, error) {
	file,head,err:=r.FormFile(key)
	if err!=nil {
		return nil,nil,err
	}
	defer file.Close()

	ext:=strings.ToLower(path.Ext(head.Filename))
	if find(ext[1:],allowedTypes)<0 {
		return nil,nil,errors.New("FILE TYPE SHOULD BE "+strings.Join(allowedTypes,","))
	}

	if int(head.Size)>allowedMaxSize {
		return nil,nil,errors.New("FILE SIZE EXCEED")
	}

	m:=make([]byte,head.Size)
	if _,err:=file.Read(m);err!=nil {
		return nil,nil,err
	}
	return m,head,nil
}
