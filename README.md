# Astral API

> [!IMPORTANT]  
> This code has been publicly released and archived. You may use it in accordance with the [License](LICENSE). **Support will not be provided.**

# Getting Started

> [!NOTE]
> Ensure [Git](https://git-scm.org) and [Go](https://go.dev) are installed.

## Clone the repository
```
git clone https://github.com/astralservices/api.git
```

## Rename / Copy the example environment file, then fill out the variables in the file
```
mv example.env .env.local
```

## Download the Go modules
```
go mod download
```

## Start the bots server
```
go run .
```

Alternatively, use the VSCode launch configuration to have debugging enabled
