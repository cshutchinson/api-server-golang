package main

import(
  "database/sql"
  "github.com/gin-gonic/gin"
  _ "github.com/lib/pq"
  "gopkg.in/gorp.v1"
  "log"
  "strconv"
)

type User struct {
  Id int64 `db:"id" json:"id"`
  Firstname string `db:"firstname" json:"firstname"`
  Lastname string `db:"lastname" json: "lastname"`
}

func main() {
  r := gin.Default()

  v1 := r.Group("api/v1")
    {
    v1.GET("/users", GetUsers)
    v1.GET("/users/:id", GetUser)
    v1.POST("/users", PostUser)
    v1.PUT("/users/:id", UpdateUser)
    v1.DELETE("/users/:id", DeleteUser)
    }

  r.Run(":8080")
}

func GetUsers(c *gin.Context){
  var users []User
  _, err := dbmap.Select(&users, "SELECT * FROM users")

  if err==nil{
    c.JSON(200, users)
  } else {
    c.JSON(404, gin.H{"error": "no user(s) in users table"})
  }
  // curl -i http://localhost:8080/api/v1/users
}

func GetUser(c *gin.Context) {
  id := c.Params.ByName("id")
  var user User
  err := dbmap.SelectOne(&user, "SELECT * FROM users WHERE id=?", id)

  if err == nil {
    user_id, _ := strconv.ParseInt(id, 0, 64)

  content := &User{
    Id: user_id,
    Firstname: user.Firstname,
    Lastname: user.Lastname,
  }
  c.JSON(200, content)
  } else {
  c.JSON(404, gin.H{"error": "user not found"})
  }
// curl -i http://localhost:8080/api/v1/users/1
}

func PostUser(c *gin.Context) {
  var user User
  c.Bind(&user)


  if user.Firstname != "" && user.Lastname != "" {
    db, _ := sql.Open("postgres", "user=cshutchinson dbname=godb sslmode=disable")
    // checkErr(err)
    var lastInsertId int64
    _ = db.QueryRow("INSERT INTO users(firstname,lastname, id) VALUES($1,$2, Default) returning id;", user.Firstname, user.Lastname).Scan(&lastInsertId)
    content := &User{
      Id: lastInsertId,
      Firstname: user.Firstname,
      Lastname: user.Lastname,
    }
    c.JSON(201, content)
    // if insert, _ := dbmap.Exec(`INSERT INTO users (firstname, lastname) VALUES (?, ?) returning id;`, user.Firstname, user.Lastname).Scan(&lastInsertId); insert != nil {
    //   user_id, err := insert.LastInsertId()
    //   if err == nil {
    //     content := &User{
    //       Id: user_id,
    //       Firstname: user.Firstname,
    //       Lastname: user.Lastname,
    //     }
    //     c.JSON(201, content)
    //   } else {
    //     checkErr(err, "Insert failed")
    //   }
    // }
  } else {
  c.JSON(422, gin.H{"error": "fields are empty"})
  }
// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func UpdateUser(c *gin.Context){

}

func DeleteUser(c *gin.Context){

}

var dbmap = initDb()
func initDb() *gorp.DbMap {
 db, err := sql.Open("postgres", "user=cshutchinson dbname=godb sslmode=disable")
 checkErr(err, "sql.Open failed")
 dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
 dbmap.AddTableWithName(User{}, "users").SetKeys(true, "id")
 err = dbmap.CreateTablesIfNotExists()
 checkErr(err, "Create table failed")
return dbmap
}
func checkErr(err error, msg string) {
 if err != nil {
 log.Fatalln(msg, err)
 }
}
