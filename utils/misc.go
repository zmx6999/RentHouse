package utils

import (
	"net/http"
	"mime/multipart"
	"path"
	"errors"
	"strings"
		"github.com/weilaihui/fdfs_client"
)

func Find(x string,a []string) int {
	n:=len(a)
	for i:=0; i<n; i++ {
		if a[i]==x {
			return i
		}
	}
	return -1
}

func PrepareUploadFile(r *http.Request,key string,allowTypes []string,allowMaxSize int) ([]byte,*multipart.FileHeader,error) {
	file,head,err:=r.FormFile(key)
	if err!=nil {
		return nil,nil,err
	}
	defer file.Close()

	ext:=strings.ToLower(path.Ext(head.Filename))
	if Find(ext[1:],allowTypes)<0 {
		return nil,nil,errors.New("FILE TYPE SHOULD BE "+strings.Join(allowTypes,","))
	}
	if head.Size>int64(allowMaxSize) {
		return nil,nil,errors.New("FILE SIZE EXCEEDS")
	}

	data:=make([]byte,head.Size)
	_,err=file.Read(data)
	if err!=nil {
		return nil,nil,err
	}

	return data,head,nil
}

func UploadFile(data []byte,ext string) (string,error) {
	client,err:=fdfs_client.NewFdfsClient(FDFSConfig)
	if err!=nil {
		return "",err
	}

	r,err:=client.UploadByBuffer(data,ext)
	if err!=nil {
		return "",err
	}

	return r.RemoteFileId,nil
}

func GetParam(key string,r *http.Request) string {
	q:=r.URL.Query()[key]
	if len(q)>0 {
		return q[0]
	} else {
		return ""
	}
}
