# Terraform Infrastructure

Infrastructure as Code for AI FinOps platform on AWS.

## Module Layout

```
modules/
├── network/          # VPC, subnets, security groups
├── eks/              # EKS cluster, node groups, IRSA
├── rds/              # PostgreSQL 16 Multi-AZ
├── elasticache/      # Redis 7
├── msk/              # Kafka (MSK)
├── s3/               # Object storage
├── kms/              # Encryption keys
├── route53/          # DNS
└── observability/    # Prometheus, Grafana (or AMP)
```

## Environments

- `environments/dev/` — minimal EKS for staging
- `environments/staging/` — pre-production
- `environments/production/` — full production stack

## State

S3 backend + DynamoDB lock table per environment.

## Usage (Phase 1+)

```bash
cd environments/staging
terraform init
terraform plan
terraform apply
```

See [docs/phase-0.5/10-infrastructure.md](../../docs/phase-0.5/10-infrastructure.md) for full spec.
