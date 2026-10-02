# SNAP Lab (Go)

A hands-on implementation of the Bank Indonesia SNAP v1.0.2 security protocol:
a mock PJP AIS server plus a SNAP consumer client, sharing one protocol package.

- Access Token B2B (service code 73)
- Balance Inquiry (service code 11)
- Transfer Credit / Intrabank (service code 17)

## Quick start
```bash
make keys      # generate RSA keypair
make server    # terminal 1
make token     # terminal 2
make balance
make transfer
```
