# IQ Theory - AWS Deployment Architecture

## AWS Services Required

### Core Infrastructure

- **Amazon ECS (Elastic Container Service)** - Container orchestration
- **Amazon RDS (PostgreSQL)** - Managed database
- **Application Load Balancer (ALB)** - Load balancing and SSL termination
- **Amazon VPC** - Network isolation
- **Amazon ECR (Elastic Container Registry)** - Docker image storage

### Storage & CDN

- **Amazon S3** - Static asset storage (note images, frontend build)
- **Amazon CloudFront** - CDN for global content delivery

### Security & Authentication

- **AWS Certificate Manager** - SSL/TLS certificates
- **AWS Secrets Manager** - Database credentials and API keys
- **AWS IAM** - Identity and access management

### Monitoring & Logging

- **Amazon CloudWatch** - Monitoring and logging
- **AWS X-Ray** - Distributed tracing (optional)

### Optional Enhancements

- **Amazon ElastiCache (Redis)** - Session storage and caching
- **Amazon SQS** - Background job processing
- **AWS Lambda** - Serverless functions for specific tasks

## Architecture Overview

```
Internet → CloudFront → ALB → ECS Tasks (Go API) → RDS (PostgreSQL)
                    ↘ S3 (Static Assets)
```

## Deployment Steps

### 1. VPC and Networking Setup

```bash
# Create VPC with public and private subnets across 2 AZs
# Public subnets: ALB
# Private subnets: ECS tasks, RDS
```

### 2. RDS PostgreSQL Setup

```bash
# Multi-AZ deployment for high availability
# Instance class: db.t3.micro (for development) or db.t3.small (production)
# Storage: 20GB SSD with auto-scaling
# Backup retention: 7 days
# Security group: Allow access only from ECS tasks
```

### 3. ECR Repository Creation

```bash
# Create repositories for:
# - iq-theory-api (Go backend)
# - iq-theory-frontend (React app with nginx)
```

### 4. ECS Cluster and Services

```bash
# Fargate launch type for serverless containers
# Auto-scaling based on CPU/memory usage
# Service discovery for internal communication
```

### 5. Application Load Balancer

```bash
# Internet-facing ALB
# SSL termination with ACM certificate
# Path-based routing:
#   - /api/* → Backend service
#   - /* → Frontend service or S3/CloudFront
```

### 6. S3 and CloudFront

```bash
# S3 bucket for:
#   - Note images (/assets/notes/)
#   - Frontend static files (if not using containerized approach)
# CloudFront distribution for global delivery
```

## Cost Estimation (US East-1, Monthly)

### Development Environment

- **RDS db.t3.micro**: ~$13/month
- **ECS Fargate** (2 tasks, 0.25 vCPU, 0.5GB): ~$6/month
- **ALB**: ~$16/month
- **S3 + CloudFront**: ~$5/month
- **Other services**: ~$10/month
- **Total**: ~$50/month

### Production Environment

- **RDS db.t3.small (Multi-AZ)**: ~$50/month
- **ECS Fargate** (4 tasks, 0.5 vCPU, 1GB): ~$25/month
- **ALB**: ~$16/month
- **S3 + CloudFront**: ~$15/month
- **ElastiCache**: ~$15/month
- **Other services**: ~$20/month
- **Total**: ~$140/month

## Docker Configuration

### Backend Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Frontend Dockerfile

```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## Environment Variables

### Backend Environment Variables

```bash
DATABASE_URL=postgresql://username:password@rds-endpoint:5432/iqtheory
REDIS_URL=redis://elasticache-endpoint:6379
JWT_SECRET=your-jwt-secret
AWS_REGION=us-east-1
S3_BUCKET_NAME=iq-theory-assets
CORS_ORIGINS=https://yourdomain.com
```

### Frontend Environment Variables

```bash
REACT_APP_API_URL=https://api.yourdomain.com
REACT_APP_CDN_URL=https://d123456789.cloudfront.net
```

## Terraform Infrastructure as Code Example

### Main Configuration

```hcl
provider "aws" {
  region = var.aws_region
}

# VPC
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "iq-theory-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["${var.aws_region}a", "${var.aws_region}b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_nat_gateway = true
  enable_vpn_gateway = false
}

# RDS
resource "aws_db_instance" "postgres" {
  identifier = "iq-theory-db"

  engine         = "postgres"
  engine_version = "15.4"
  instance_class = "db.t3.micro"

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_encrypted     = true

  db_name  = "iqtheory"
  username = "iqtheory_user"
  password = var.db_password

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.default.name

  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"

  skip_final_snapshot = false
  final_snapshot_identifier = "iq-theory-final-snapshot"
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "iq-theory-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}
```

## CI/CD Pipeline with GitHub Actions

```yaml
name: Deploy to AWS

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      - name: Build and push backend image
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: iq-theory-api
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG ./server
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG

      - name: Deploy to ECS
        run: |
          aws ecs update-service --cluster iq-theory-cluster --service iq-theory-api --force-new-deployment
```

## Security Best Practices

1. **Network Security**

   - Use private subnets for ECS tasks and RDS
   - Security groups with least privilege access
   - VPC endpoints for AWS services

2. **Application Security**

   - Store secrets in AWS Secrets Manager
   - Use IAM roles for ECS tasks
   - Enable encryption at rest and in transit

3. **Database Security**

   - Enable SSL connections
   - Regular automated backups
   - Database activity monitoring

4. **Monitoring and Alerting**
   - CloudWatch alarms for key metrics
   - Log aggregation and analysis
   - Performance monitoring

This architecture provides a scalable, secure, and cost-effective deployment for IQ Theory on AWS using Docker containers.
