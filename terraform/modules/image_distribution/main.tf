
locals {
  s3_origin_id  = "images-s3-origin"
}

resource "aws_s3_bucket" "image_bucket" {
  bucket = "bluthinator-images"

  tags = {
    Name = "Bluthinator Image Bucket"
  }
}

resource "aws_s3_bucket_acl" "image_bucket_acl" {
  bucket = aws_s3_bucket.image_bucket.bucket

  acl = "private"
}

resource "aws_cloudfront_origin_access_identity" "image_bucket_origin_access_identity" {
  comment = "Access to Bluthinator Image Bucket"
}

resource "aws_s3_bucket_policy" "image_bucket_policy" {
  bucket = aws_s3_bucket.image_bucket.bucket

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Action    = "s3:GetObject",
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.image_bucket_origin_access_identity.iam_arn
        },
        Resource  = "${aws_s3_bucket.image_bucket.arn}/*"
      }
    ]
  })
}

resource "aws_cloudfront_distribution" "image_distribution" {
  origin {
    domain_name = aws_s3_bucket.image_bucket.bucket_regional_domain_name
    origin_id   = local.s3_origin_id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.image_bucket_origin_access_identity.cloudfront_access_identity_path
    }
  }

  enabled = true
  is_ipv6_enabled = true
  aliases = [
    "img.bluthinator.com"
  ]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA", "GB", "DE"]
    }
  }
  
  viewer_certificate {
    acm_certificate_arn = var.acm_certificate_arn
    cloudfront_default_certificate = false
    minimum_protocol_version = "TLSv1.2_2021"
    ssl_support_method = "sni-only"
  }
}