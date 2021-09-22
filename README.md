# keyv
Keyv provides a key-value interface to access database. By changing the adapter, you can change the database you use.  
Inspired by [lukechilds/keyv](https://github.com/lukechilds/keyv), a similar nodejs module.  

# Adapters
[sqlite3](https://github.com/simba-fs/keyvSqlite3)

# Usage
```go
import _ "github.com/simba-fs/keyvSqlite3"
import "github.com/simba-fs/keyv"

func main(){
	db, err := keyv.New()

}

```
