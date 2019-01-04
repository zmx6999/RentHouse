package utils

import "github.com/weilaihui/fdfs_client"

func UploadFile(fileBytes []byte,fileExt string) (string,error) {
	client,err:=fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	if err!=nil {
		return "",err
	}
	re,err:=client.UploadByBuffer(fileBytes,fileExt)
	if err!=nil {
		return "",err
	}
	return re.RemoteFileId,nil
}
