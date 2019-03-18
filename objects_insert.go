package gogcs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"google.golang.org/api/storage/v1"

	"github.com/gorilla/mux"

	bolt "go.etcd.io/bbolt"
)

func (s *Server) objectInsert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	switch r.URL.Query().Get("uploadType") {
	case "media":
		{
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Print(err)
				http.Error(w, "couldn't read the body", http.StatusBadRequest)
				return
			}

			oName := r.URL.Query().Get("name")
			if oName == "" {
				http.Error(w, "couldn't get object name", http.StatusBadRequest)
				return
			}

			err = s.db.Update(func(tx *bolt.Tx) error {
				bName, err := getBucketName(bucket, tx)
				if err != nil {
					log.Print(err)
					return err
				}

				err = createObject(string(bName), oName, body)
				if err != nil {
					log.Print(err)
					return err
				}

				return tx.Bucket(bName).Put([]byte(oName), getBytesFromObject(storage.Object{
					Name: oName,
				}))
			})
			if err != nil {
				http.Error(w, "couldn't update the file system and database", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "%s", string(getBytesFromObject(storage.Object{
				Name: oName,
			})))
		}
	case "multipart":
		{
			mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil {
				log.Print(err)
				http.Error(w, "couldn't parse media", http.StatusBadRequest)
				return
			}
			if strings.HasPrefix(mediaType, "multipart/") {
				mr := multipart.NewReader(r.Body, params["boundary"])
				for {
					p, err := mr.NextPart()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Print(err)
						http.Error(w, "couldn't get next part in multipart content", http.StatusBadRequest)
						return
					}
					slurp, err := ioutil.ReadAll(p)
					if err != nil {
						log.Print(err)
						http.Error(w, "couldn't read the next part", http.StatusBadRequest)
						return
					}
					fmt.Printf("Part: %q\n", slurp)
				}
			}
		}
	case "resumable":
		{

		}
	}

	/*r.FormF
	bucket := ""
	err = s.db.Update(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			object := storage.Object{}
			if strings.Contains(string(name), bucket) {
				b.Put([]byte(object.Name), nil)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		http.Error(w, "couldn't update database", http.StatusBadRequest)
		return
	}*/
}
