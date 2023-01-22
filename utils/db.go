package utils

import(
	"database/sql"
    "fmt"
    "context"
    "github.com/go-redis/redis/v8"
    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "1234"
    dbname   = "shaparak"
)

func ConnectToDb() {
    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
        
   db, err := sql.Open("postgres", psqlconn)
   CheckError(err)
	
   defer db.Close()

   err = db.Ping()
   CheckError(err)

   fmt.Println("Connected!")
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func GetUserData(username string, password string) {
    ctx := context.Background()
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", 
        DB:       0,  
    })

    val, err := client.Get(ctx, username).Result()
    
    CheckError(err)
    if (val != "") {
        fmt.Println(val)
    } else {
        fmt.Println("key not found")
    }
}