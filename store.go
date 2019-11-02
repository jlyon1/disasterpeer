package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/asdine/storm"
	"github.com/google/uuid"
)

type Store struct {
	db *storm.DB
}

func NewStore() (*Store, error) {
	db, err := storm.Open("disaster.db")
	if err != nil {
		fmt.Println("open error")
		return nil, err
	}

	// init buckets for each struct
	for _, t := range []interface{}{
		&MyInfo{},
		&EncryptedMessage{},
	} {
		if err := db.Init(t); err != nil {
			fmt.Println("can't create bucket")
			return nil, err
		}
	}

	s := &Store{
		db: db,
	}

	return s, nil
}

// type MetaImages struct {
// 	Image []byte
// }

type MyInfo struct {
	ID    uuid.UUID // user ID
	Name  string
	Email string
	Phone string
	Lat   float64
	Long  float64
	Time  time.Time
	// Meta  MetaImages
}

func (s *Store) SetMyInfo(info *MyInfo) error {
	return s.db.Save(info)
}

func (s *Store) UpdateLocation(myID uuid.UUID, newLat float64, newLong float64) error {
	return s.db.Update(&MyInfo{ID: myID, Lat: newLat, Long: newLong, Time: time.Now()})
}

type EncryptedMessage struct {
	ID   int `storm:"id,increment"` // message ID
	Sent time.Time
	Body []byte
}

// Throw away old id from
func (s *Store) SaveMessages(msg []byte) {
	bytes := []byte(msg)
	var messages []EncryptedMessage
	err := json.Unmarshal(bytes, &messages)
	if err != nil {
		log.Panicln("couldn't unmarshal message")
	}

	for _, m := range messages {
		if err = s.db.Save(&m); err != nil {
			log.Panicln("message error")
		}
	}
}

func (s *Store) GetAllMessages(myID uuid.UUID) []byte {
	// Get all existing messages
	var messages []EncryptedMessage
	err := s.db.All(&messages)

	// Append user's most recent info onto encrypted messages
	rng := rand.Reader

	var myInfo MyInfo
	bodyString, err := json.Marshal(s.db.One("UUID", myID, &myInfo))

	pub, err := ioutil.ReadFile("key.pem")
	if err != nil {
		fmt.Println(err)
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		fmt.Println("bad")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)

	var pubKey *rsa.PublicKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		log.Panicln("Unable to parse RSA public key, generating a temp one")
	}

	body, _ := rsa.EncryptOAEP(sha256.New(), rng, pubKey, bodyString, []byte(myID.String()))
	newMessage := EncryptedMessage{
		Sent: time.Now(),
		Body: body,
	}

	messages = append(messages, newMessage)

	// Marshal messages
	bytes, err := json.Marshal(messages)
	if err != nil {
		log.Panicln("couldn't marshal messages")
	}

	return bytes
}
