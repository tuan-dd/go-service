env "default" {
  host = getenv("HOST")
  port = getenv("PORT")
  user_name = getenv("USERNAME")
  password = getenv("PASSWORD")
  dbName = getenv("DBNAME")
  sslMode   = getenv("SSL_MODE")

  migration {
    dir = "file://migrations"
  }
}
