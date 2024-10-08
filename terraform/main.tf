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

resource "aws_security_group" "bastion_server_sg" {
  vpc_id = module.vpc.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [ var.basion_server_ingress_cidr ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "bastion_server" {
  ami           = "ami-0e86e20dae9224db8" # Ubuntu 22.04 LTS (x86_64)
  instance_type = "t2.small"
  key_name      = "bluthinator-bastion-server-key"
  subnet_id     = module.vpc.public_subnet_ids[0]
  vpc_security_group_ids = [aws_security_group.bastion_server_sg.id]

  tags = {
    Name = "Bluthinator Bastion Server"
  }
}

resource "aws_instance" "elastic_search_server" {
  ami           = "ami-096ea6a12ea24a797" # Ubuntu 22.04 LTS (arm)
  instance_type = "c6g.medium"
  key_name      = "bluthinator-elasticsearch-server-key"
  subnet_id     = module.vpc.private_subnet_ids[0]
  vpc_security_group_ids = [aws_security_group.elasticsearch_sg.id]

  tags = {
    Name = "Bluthinator ElasticSearch Server"
  }
}

resource "aws_security_group" "elasticsearch_sg" {
  name        = "elasticsearch-sg"
  description = "Allow access to Elasticsearch from the bastion server"
  vpc_id      = module.vpc.vpc_id

  tags = {
    Name = "Elasticsearch Security Group"
  }
}

resource "aws_vpc_security_group_ingress_rule" "elasticsearch_sg_ingress_rule_es" {
  security_group_id = aws_security_group.elasticsearch_sg.id
  from_port         = 9200
  to_port           = 9200
  ip_protocol       = "tcp"
  referenced_security_group_id = aws_security_group.bastion_server_sg.id
}

resource "aws_vpc_security_group_ingress_rule" "elasticsearch_sg_ingress_rule_ssh" {
  security_group_id = aws_security_group.elasticsearch_sg.id
  from_port         = 22
  to_port           = 22
  ip_protocol       = "tcp"
  referenced_security_group_id = aws_security_group.bastion_server_sg.id
}

resource "aws_vpc_security_group_egress_rule" "elasticsearch_sg_egress_rule" {
  security_group_id = aws_security_group.elasticsearch_sg.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
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

resource "aws_ecr_repository" "api_image_repo" {
  name = "bluthinator-api"

  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "Bluthinator API Image Repository"
  }
}

resource "aws_ecr_repository" "web_image_repo" {
  name = "bluthinator-web"

  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "Bluthinator Web Server Image Repository"
  }
}