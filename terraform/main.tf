terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      awsApplication = var.aws_application_rg
    }
  }
}

locals {
  acm_certificate_arn = "arn:aws:acm:us-east-1:635676917059:certificate/a45a94db-250c-43d7-b463-432c2617c251"
}

module "vpc" {
  source = "./modules/vpc"

  cidr_block           = var.vpc_cidr_block
  azs                  = var.azs
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
}

module "rds" {
  source = "./modules/rds"

  username           = var.db_username
  password           = var.db_password
  security_group_ids = [aws_security_group.db_security_group.id]
  subnet_group_name  = aws_db_subnet_group.database_subnet_group.name
}

module "image_distribution" {
  source = "./modules/image_distribution"

  s3_bucket_name = "bluthinator-images"
  acm_certificate_arn = local.acm_certificate_arn
}

resource "aws_db_subnet_group" "database_subnet_group" {
  name       = "database-subnet-group"
  subnet_ids = module.vpc.private_subnet_ids

  tags = {
    Name = "Bluthinator Database Subnet Group"
  }
}

resource "aws_security_group" "db_security_group" {
  vpc_id = module.vpc.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "Bluthinator DB Security Group"
  }
}