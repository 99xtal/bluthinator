variable "acm_certificate_arn" {
  description = "The ARN of the ACM certificate to use for HTTPS"
  type = string
}

variable "s3_bucket_name" {
  description = "The name of the S3 bucket to use for the CloudFront distribution"
  type = string
}