resource "aws_ecr_repository" "app" {
  name                 = "${var.project_name}-repo"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_dynamodb_table" "url_table" {
  name         = "${var.project_name}-urls"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ShortCode"

  attribute {
    name = "ShortCode"
    type = "S"
  }
}

resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-cluster"
}