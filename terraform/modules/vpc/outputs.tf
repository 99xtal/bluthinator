output "private_subnet_ids" {
  value = aws_subnet.private_subnets[*].id
}

output "vpc_id" {
  value = aws_vpc.vpc.id
}