package keyv

import (
	"encoding/json"
	"errors"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type AdapterNewer interface {
	Connect(uri string) (Adapter, error)
}

type Adapter interface {
	Has(key string) bool
	Get(key string) (string, error)
	Set(key string, val string) error
	Remove(key string) error
	Keys() []string
}

type Keyv struct {
	AdapterName string
	Adapter     Adapter
	Uri         string
	Namespace   string
}

// addNS add namespace to key
func (k *Keyv) addNS(key string) string {
	return "keyv:" + k.Namespace + ":" + key
}

// Get value from DB with key
func (k *Keyv) Get(key string, v interface{}) error {
	// add namespace
	key = k.addNS(key)

	// check if key exists
	if !k.Adapter.Has(key) {
		return nil
	}

	// get raw data
	data, err := k.Adapter.Get(key)
	if err != nil {
		return err
	}

	// convert json into go struct
	return json.Unmarshal([]byte(data), v)
}

// Set value with key in DB
func (k *Keyv) Set(key string, value interface{}) error {
	// add namespace
	key = k.addNS(key)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return k.Adapter.Set(key, string(data))
}

// Remove data with key from DB
func (k *Keyv) Remove(key string) error {
	// add namespace
	key = k.addNS(key)

	return k.Remove(key)
}

// Clear remove all data in DB with the same namespace
func (k *Keyv) Clear() error {
	keys := k.Adapter.Keys()
	for _, key := range keys {
		if err := k.Adapter.Remove(key); err != nil {
			return err
		}
	}

	return nil
}

var (
	ErrAdapterNewerNotFound   = errors.New("adapter newer not found")
	ErrAdapterNewerNameExists = errors.New("adapter newer name exists")
)

var adapterNewers = map[string]AdapterNewer{}

// Register add a new adapter newer with name
func Register(name string, adapterNewer AdapterNewer) error {
	_, ok := adapterNewers[name]
	if ok {
		return ErrAdapterNewerNameExists
	}

	adapterNewers[name] = adapterNewer
	return nil
}

// New create a keyv object
func New(uri string, namespace string) (*Keyv, error) {
	adapterName := strings.SplitN(uri, "://", 2)[0]

	newer, ok := adapterNewers[adapterName]
	if !ok {
		return nil, ErrAdapterNewerNotFound
	}

	adapter, err := newer.Connect(uri)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = "default"
	}

	keyv := &Keyv{
		AdapterName: adapterName,
		Adapter:     adapter,
		Uri:         uri,
		Namespace:   namespace,
	}

	return keyv, nil
}
