package gogcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	bolt "github.com/boltdb/bolt"
	storage "google.golang.org/api/storage/v1"
)

func (s *Server) bucketInsert(w http.ResponseWriter, r *http.Request) {
	project := getProjectFromRequest(r)

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't get body", http.StatusBadRequest)
		return
	}

	var b storage.Bucket
	err = json.Unmarshal(bs, &b)
	if err != nil {
		http.Error(w, "couldn't unmarshal body", http.StatusBadRequest)
		return
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(getBucketConfigPathBs(project, b.Name))
		if err != nil {
			log.Println("create bucket: ", err)
			return err
		}
		log.Println("json from object: ", getBytesFromObject(b))

		return tx.Bucket(rootBucket).Put(getBucketConfigPathBs(project, b.Name), getBytesFromObject(b))
	})
	if err != nil {
		http.Error(w, "couldn't update database", http.StatusBadRequest)
		return
	}

	err = createProjectAndBucket(project, b.Name)
	if err != nil {
		http.Error(w, "couldn't create directories", http.StatusBadRequest)
		return
	}

	bs, err = json.Marshal(b)
	if err != nil {
		http.Error(w, "couldn't marshal response", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s", string(bs))
}
