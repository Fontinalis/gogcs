package gogcs

import (
	"encoding/json"
	"fmt"
	"net/http"

	bolt "go.etcd.io/bbolt"
	storage "google.golang.org/api/storage/v1"
)

func (s *Server) bucketsList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AUTHORIZATION", r.Header.Get("Authorization"))
	buckets := make([]*storage.Bucket, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucket := storage.Bucket{}
			fmt.Println(string(name))
			if string(name) == string(rootBucket) {
				return nil
			}
			bs := tx.Bucket(rootBucket).Get(name)
			fmt.Println("bucket info: ", string(bs))
			err := json.Unmarshal(bs, &bucket)
			if err != nil {
				return err
			}
			buckets = append(buckets, &bucket)
			return nil
		})
	})
	if err != nil {
		http.Error(w, "couldn't get buckets", http.StatusBadRequest)
		return
	}

	bs, err := json.Marshal(newBucketListResult(buckets))
	if err != nil {
		http.Error(w, "couldn't marshal result", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(bs))
}

type bucketListResult struct {
	Kind          string            `json:"kind"`
	NextPageToken string            `json:"nextPageToken"`
	Items         []*storage.Bucket `json:"items"`
}

func newBucketListResult(bs []*storage.Bucket) bucketListResult {
	return bucketListResult{
		Kind:  "storage#buckets",
		Items: bs,
	}
}
