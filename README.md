# driftctl-diff

> A lightweight utility that compares Terraform state files across environments and surfaces configuration drift in a readable report.

---

## Installation

```bash
go install github.com/yourusername/driftctl-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftctl-diff.git
cd driftctl-diff
go build -o driftctl-diff .
```

---

## Usage

Compare two Terraform state files and generate a drift report:

```bash
driftctl-diff --base staging.tfstate --target production.tfstate
```

Output to a file:

```bash
driftctl-diff --base staging.tfstate --target production.tfstate --output report.txt
```

**Example output:**

```
[DRIFT DETECTED]
  ~ aws_instance.web_server
      instance_type: "t2.micro" → "t3.medium"
      ami:           "ami-0abcdef" → "ami-0fedcba"

  + aws_s3_bucket.logs  (present in production, missing in staging)

Summary: 2 drifted resources, 1 missing resource
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Path to the base state file | required |
| `--target` | Path to the target state file | required |
| `--output` | Write report to a file instead of stdout | — |
| `--format` | Output format: `text`, `json` | `text` |

---

## License

This project is licensed under the [MIT License](LICENSE).