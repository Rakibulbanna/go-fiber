variable "db_url" {
  type = string
  default = "postgres://postgres:postgres@localhost:5432/fiber_demo?sslmode=disable"
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/atlas/main.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  url = var.db_url
  dev = "docker://postgres/15/dev?search_path=public"
  migration {
    dir = "file://migrations"
  }
}

env "local" {
  src = "file://migrations"
  url = var.db_url
  dev = "docker://postgres/15/dev?search_path=public"
}

