# Inventory

Is an open source Inventory API

## Get Started
- git clone git@github.com:jacky-htg/inventory.git
- cp .env.example .env
- edit .env with your environment
- go mod init github.com/jacky-htg/inventory
- go run cmd/main.go migrate
- go run cmd/main.go seed
- go run cmd/main.go scan-access
- go test (for run this command, you need docker installed in your laptop)
- go run main.go