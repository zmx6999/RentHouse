package utils

import (
	"net/http"
	"mime/multipart"
	"strings"
	"path"
	"errors"
	"github.com/weilaihui/fdfs_client"
	"crypto/sha512"
	"encoding/hex"
	"github.com/garyburd/redigo/redis"
)

const (
	ORDER_STATUS_WAIT_ACCEPT= "WAIT_ACCEPT"
	ORDER_STATUS_WAIT_PAYMENT= "WAIT_PAYMENT"
	ORDER_STATUS_PAID= "PAID"
	ORDER_STATUS_WAIT_COMMENT= "WAIT_COMMENT"
	ORDER_STATUS_COMPLETE= "COMPLETE"
	ORDER_STATUS_CANCELED= "CONCELED"
	ORDER_STATUS_REJECTED= "REJECTED"
)

func AddDomain2Url(url string) string {
	return "http://"+FDFSHost+":"+FDFSPort+"/"+url
}

func PrepareUpload(r *http.Request,key string,allowedTypes []string,allowedMaxSize int64) ([]byte,*multipart.FileHeader,error) {
	file,head,err:=r.FormFile(key)
	if err!=nil {
		return nil,nil,err
	}
	defer file.Close()

	ext:=strings.ToLower(path.Ext(head.Filename))
	if Find(allowedTypes,ext[1:])<0 {
		return nil,nil,errors.New("FILE TYPES CAN ONLY BE "+strings.Join(allowedTypes,","))
	}
	if head.Size>allowedMaxSize {
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
	client,err:=fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	if err!=nil {
		return "",err
	}
	re,err:=client.UploadByBuffer(data,ext)
	if err!=nil {
		return "",err
	}
	return re.RemoteFileId,nil
}

func Find(a []string,x string) int {
	for i:=0; i<len(a); i++ {
		if a[i]==x {
			return i
		}
	}
	return -1
}

func GetParam(key string,r *http.Request) string {
	q:=r.URL.Query()
	if len(q[key])>0 {
		return q[key][0]
	} else {
		return ""
	}
}

func Sha512Str(x string) string {
	h:=sha512.New()
	h.Write([]byte(x))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUserId(sessionId string) (int,error) {
	conn,err:=redis.Dial("tcp",RedisHost+":"+RedisPort)
	if err!=nil {
		return 0,err
	}
	defer conn.Close()

	userId,err:=redis.Int(conn.Do("get",sessionId+"_user_id"))
	if err!=nil {
		return 0,err
	}
	return userId,nil
}
