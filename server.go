package gogcs

import (
	"net/http"

	bolt "go.etcd.io/bbolt"
	"github.com/gorilla/mux"
)

var rootBucket = []byte(".")

func New() (*Server, error) {
	s := &Server{}
	s.r = mux.NewRouter()
	// Buckets
	s.r.Path("/b").Methods("POST").HandlerFunc(s.bucketInsert) // Buckets: insert
	s.r.Path("/b").Methods("GET").HandlerFunc(s.bucketsList)   // Buckets: list

	// Objects
	s.r.Path("/b/{bucket}/o").Methods("POST").HandlerFunc(s.objectInsert) // Objects: insert

	db, err := bolt.Open("gcs.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	s.db = db
	s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(rootBucket)
		if err != nil {
			return err
		}
		return nil
	})
	return s, nil
}

type Server struct {
	r  *mux.Router
	db *bolt.DB
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
