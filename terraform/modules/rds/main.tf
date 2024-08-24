resource "aws_db_instance" "database" {
  allocated_storage = 10
  identifier        = "bluthinator-db"
  engine            = "postgres"
  instance_class    = "db.t3.micro"
  username          = var.username
  password          = var.password

  vpc_security_group_ids = var.security_group_ids
  db_subnet_group_name   = var.subnet_group_name

  skip_final_snapshot = true

  monitoring_interval = 60
  monitoring_role_arn = aws_iam_role.rds_monitoring_role.arn

  performance_insights_enabled = true

  tags = {
    Name = "Bluthinator Database"
  }
}

resource "aws_iam_role" "rds_monitoring_role" {
  name = "rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "monitoring.rds.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "Bluthinator RDS Monitoring Role"
  }
}

resource "aws_iam_policy_attachment" "rds_monitoring_attachment" {
  name       = "rds-monitoring-attachment"
  roles      = [aws_iam_role.rds_monitoring_role.name]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

