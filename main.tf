terraform {
  required_version = ">= 1.0.0"

  backend "s3" {
    bucket         = "finops-scanner-tfstate-achh-914"
    key            = "global/s3/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "finops-scanner-tflocks"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "terraform_state" {
  bucket        = "finops-scanner-tfstate-achh-914"
  force_destroy = true
}

resource "aws_dynamodb_table" "terraform_locks" {
  name         = "finops-scanner-tflocks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
