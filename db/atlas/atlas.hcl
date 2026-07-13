env "local" {
  url = getenv("DATABASE_URL")
  dev = getenv("ATLAS_DEV_URL")

  migration {
    dir = "file://migrations"
  }

  schema {
    src = "file://schemas/postgres.sql"
  }
}

env "staging" {
  url = getenv("DATABASE_URL")
  migration { dir = "file://migrations" }
}

env "production" {
  url = getenv("DATABASE_URL")
  migration { dir = "file://migrations" }
}
