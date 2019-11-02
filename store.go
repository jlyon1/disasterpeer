package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
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

func NewStore(name string) (*Store, error) {
	db, err := storm.Open(name)
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
	ID     int       `storm:"id,increment"`
	UserID uuid.UUID // user ID
	Name   string
	Email  string
	Phone  string
	Lat    float64
	Long   float64
	Time   time.Time
	// Meta  MetaImages
}

func (s *Store) SetMyInfo(newInfo *MyInfo) error {

	// Append user's most recent info onto encrypted messages
	rng := rand.Reader
	pub, err := ioutil.ReadFile("mykey.pub")
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
		fmt.Println("Unable to parse RSA public key")
	}

	bodyString, err := json.Marshal(&newInfo)

	if err != nil {
		fmt.Println(err)
	}
	body, err := rsa.EncryptOAEP(sha512.New(), rng, pubKey, bodyString, nil)
	if err != nil {
		fmt.Println(err)
	}

	newMessage := EncryptedMessage{
		Sent: time.Now(),
		Body: body,
	}
	s.db.Save(&newMessage)
	return s.db.Save(newInfo)
}

func (s *Store) GetMyInfo(myId uuid.UUID) (myInfo MyInfo, myErr error) {
	var tmp []MyInfo
	err := s.db.Find("UserID", myId, &tmp, storm.Reverse())
	return tmp[0], err
}

func (s *Store) UpdateLocation(myID uuid.UUID, newLat float64, newLong float64) error {
	myInfo, err := s.GetMyInfo(myID)
	if err != nil {
		fmt.Println(err)
	}

	newInfo := &MyInfo{
		UserID: myInfo.UserID,
		Name:   myInfo.Name,
		Email:  myInfo.Email,
		Phone:  myInfo.Phone,
		Lat:    newLat,
		Long:   newLong,
		Time:   time.Now(),
	}

	return s.SetMyInfo(newInfo)
}

type EncryptedMessage struct {
	ID   int `storm:"id,increment"` // message ID
	Sent time.Time
	Body []byte
}

// TODO: Throw away old id from
func (s *Store) SaveMessages(messages []EncryptedMessage) {
	for _, m := range messages {
		m.ID = 0
		if err := s.db.Save(&m); err != nil {
			log.Panicln("message error")
		}
	}
}

// func (s *Store) ViewMessages() {
// 	var messages []EncryptedMessage
// 	err := s.db.All(&messages)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(messages)
// 	fmt.Println(len(messages))

// 	// for _, m := range messages {
// 	// 	fmt.Println(m.ID)
// 	// 	fmt.Println(m.Body) + "\n")
// 	// }

// }

func (s *Store) GetAllMessages(myID uuid.UUID) []EncryptedMessage {

	// Get all existing messages
	var messages []EncryptedMessage
	s.db.All(&messages)
	return messages
}
