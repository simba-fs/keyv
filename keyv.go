package keyv

import (
	"encoding/json"
	"errors"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrKeyNotFound = errors.New("key not found in namespace")
)

type AdapterNewer interface {
	// Cnnect return a keyv adapter by the given uri
	Connect(uri string) (Adapter, error)
}

type Adapter interface {
	// Has checks if key exists
	Has(key string) bool
	// Get returns value by key
	Get(key string) (string, error)
	// Set sets value by key
	Set(key string, val string) error
	// Remove removes value by key
	Remove(key string) error
	// Clear remove all data in this namespace
	Clear(prefix string) error
	// Keys return all keys in this namespace
	Keys() ([]string, error)
}

type Keyv struct {
	// AdapterName is the name of used adapter
	AdapterName string
	// Adapter is the used adapter
	Adapter Adapter
	// Uri is the used uri
	Uri string
	// Namespace will be automatically add before key
	Namespace string
}

// addNS add namespace to key
func (k *Keyv) addNS(key string) string {
	return "keyv:" + k.Namespace + ":" + key
}

// Has check if key exists in the db
func (k *Keyv) Has(key string) bool {
	// add namespace
	key = k.addNS(key)

	return k.Adapter.Has(key)
}

// Keys return all keys in this namespace
func (k *Keyv) Keys() ([]string, error) {
	names, err := k.Adapter.Keys()
	if err != nil {
		return names, err
	}

	results :=  make([]string, 0)

	for _, v := range(names){
		if strings.HasPrefix(v, k.Namespace) {
			results = append(results, v)
		}
	}

	return results, nil
}

// Get value from DB with key
func (k *Keyv) Get(key string, v interface{}) error {
	// add namespace
	key = k.addNS(key)

	// check if key exists
	if !k.Adapter.Has(key) {
		return ErrKeyNotFound
	}

	// get raw data
	data, err := k.Adapter.Get(key)
	if err != nil {
		return err
	}

	// convert json into go struct
	return json.Unmarshal([]byte(data), v)
}

func (k *Keyv) GetString(key string) (string, error){
	s := ""
	err := k.Get(key, &s)
	return s, err
}

func (k *Keyv) GetInt(key string) (int, error){
	s := 0
	err := k.Get(key, &s)
	return s, err
}

func (k *Keyv) GetBool(key string) (bool, error){
	s := false
	err := k.Get(key, &s)
	return s, err
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

	return k.Adapter.Remove(key)
}

// Clear remove all data in DB with the same namespace
func (k *Keyv) Clear() error {
	return k.Adapter.Clear("keyv:" + k.Namespace)
}

var (
	ErrAdapterNewerNotFound   = errors.New("adapter newer not found")
	ErrAdapterNewerNameExists = errors.New("adapter newer name exists")
)

var adapterNewers = map[string]AdapterNewer{}

// Register add a new adapter newer with name. It should only be called in adapter package
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
