variable "username" {
  description = "The username for the database."
  type        = string
}

variable "password" {
  description = "The password for the database."
  type        = string
}

variable "subnet_group_name" {
  description = "The name of the database subnet group."
  type        = string
}

variable "security_group_ids" {
  description = "The security group IDs for the database."
  type        = list(string)
}