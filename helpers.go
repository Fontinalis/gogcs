package gogcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	bolt "github.com/boltdb/bolt"
)

func getProjectFromRequest(r *http.Request) string {
	return r.URL.Query().Get("project")
}

func getProjectFromRequestBs(r *http.Request) []byte {
	return []byte(r.URL.Query().Get("project"))
}

func getBucketConfigPath(project, bucket string) string {
	return project + "/" + bucket
}

func getBucketConfigPathBs(project string, bucket string) []byte {
	return []byte(project + "/" + bucket)
}

func getBytesFromObject(x interface{}) []byte {
	bs, _ := json.Marshal(x)
	fmt.Println("json from object:", string(bs))
	return bs
}

func createProjectAndBucket(project string, bucket string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir+"/"+getBucketConfigPath(project, bucket)+"/", os.ModePerm)
}

func getBucketName(bucket string, tx *bolt.Tx) ([]byte, error) {
	var bName []byte
	err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
		if strings.Contains(string(name), "/"+bucket) {
			bName = name
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else if bName == nil {
		return nil, errors.New("Bucket not found")
	}
	return bName, nil
}

func createObject(bName string, name string, bs []byte) error {
	sv := strings.Split(name, "/")
	rName := sv[len(sv)-1]
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	oPath := dir + "/" + bName + "/"
	if strings.TrimSuffix(name, "/"+rName) != name {
		oPath += strings.TrimSuffix(name, "/"+rName)
	}
	err = os.MkdirAll(oPath, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(oPath+rName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.Write(bs)
	if err != nil {
		return err
	}
	return f.Close()
}
