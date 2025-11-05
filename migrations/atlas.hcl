variable "db_url" {
  type = string
  default = "postgres://postgres:postgres@localhost:5432/fiber_demo?sslmode=disable"
}

env "local" {
  src = "file://migrations"
  url = var.db_url
  dev = "docker://postgres/15/dev?search_path=public"
}

