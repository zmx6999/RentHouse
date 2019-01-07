package utils

import (
	"net/http"
	"crypto/sha512"
	"encoding/hex"
	"github.com/garyburd/redigo/redis"
	"github.com/weilaihui/fdfs_client"
	"mime/multipart"
	"path"
	"github.com/pkg/errors"
	"strings"
)

func AddDomain2Url(url string) string {
	return "http://"+FastDFSHost+":"+FastDFSPort+"/"+url
}

const (
	ORDER_STATUS_WAIT_ACCEPT= "WAIT_ACCEPT"
	ORDER_STATUS_WAIT_PAYMENT= "WAIT_PAYMENT"
	ORDER_STATUS_PAID= "PAID"
	ORDER_STATUS_WAIT_COMMENT= "WAIT_COMMENT"
	ORDER_STATUS_COMPLETE= "COMPLETE"
	ORDER_STATUS_CANCELED= "CONCELED"
	ORDER_STATUS_REJECTED= "REJECTED"
)

func GetParam(r *http.Request, key string) string {
	value:=""
	values:=r.URL.Query()[key]
	if len(values)>0 {
		value=values[0]
	}
	return value
}

func GetSha512Str(x string) string {
	h:=sha512.New()
	h.Write([]byte(x))
	return hex.EncodeToString(h.Sum(nil))
}

func GetSession(sessionId string) (map[string]interface{},error) {
	conn,err:=redis.Dial("tcp",RedisHost+":"+RedisPort)
	if err!=nil {
		return nil,err
	}
	defer conn.Close()
	userId,err:=redis.Int(conn.Do("get",sessionId+"_user_id"))
	if err!=nil {
		return nil,err
	}
	mobile,err:=redis.String(conn.Do("get",sessionId+"_mobile"))
	if err!=nil {
		return nil,err
	}
	name,err:=redis.String(conn.Do("get",sessionId+"_name"))
	if err!=nil {
		return nil,err
	}
	return map[string]interface{}{
		"user_id":userId,
		"mobile":mobile,
		"name":name,
	},nil
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

func PrepareUpload(r *http.Request,key string,allowedTypes []string,allowedMaxSize int64) ([]byte,*multipart.FileHeader,error) {
	file,head,err:=r.FormFile(key)
	if err!=nil {
		return nil,nil,err
	}
	defer file.Close()

	ext:=strings.ToLower(path.Ext(head.Filename))
	if Find(ext[1:],allowedTypes)<0 {
		return nil,nil,errors.New("FILE TYPE SHOULD BE "+strings.Join(allowedTypes,","))
	}
	if head.Size>allowedMaxSize {
		return nil,nil,errors.New("FILE SIZE EXCEEDS")
	}

	data:=make([]byte,int(head.Size))
	_,err=file.Read(data)
	if err!=nil {
		return nil,nil,err
	}
	return data,head,nil
}

func Find(x string,a []string) int {
	for k,v:=range a{
		if v==x {
			return k
		}
	}
	return -1
}
