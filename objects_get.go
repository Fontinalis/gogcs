package gogcs

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	bolt "github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func (s *Server) objectGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]

	switch r.URL.Query().Get("alt") {
	case "media":
		{
			var path string
			err := s.db.Update(func(tx *bolt.Tx) error {
				bn, err := getBucketName(bucket, tx)
				if err != nil {
					return err
				}
				path = string(bn)
				return nil
			})
			if err != nil {
				http.Error(w, "couldn't get bucket path", http.StatusBadRequest)
				return
			}
			dir, err := os.Getwd()
			if err != nil {
				http.Error(w, "couldn't get working dir", http.StatusBadRequest)
				return
			}
			log.Println("path: ", dir+"/"+path+"/"+object)
			f, err := os.Open(dir + "/" + path + "/" + object)
			if err != nil {
				http.Error(w, "couldn't open object", http.StatusBadRequest)
				return
			}
			bs, err := ioutil.ReadAll(f)
			if err != nil {
				http.Error(w, "couldn't read object", http.StatusBadRequest)
				return
			}
			log.Println(string(bs))
			w.Write(bs)
		}
	default:
		{
		}
	}
}
