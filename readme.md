# CountrySearch API

A simple Go-based REST API that fetches country information (name, capital, currency symbol, and population) using the [REST Countries API](https://restcountries.com/) with in-memory **LRU caching and TTL** support.

---

##  Features

-  Search countries by name
-  In-memory LRU cache with TTL for performance
-  Timeout-safe API calls using `context`
-  100% test coverage (unit tests)
-  Clean and idiomatic Go project structure

---

##  Setup Instructions

###  Prerequisites
- Go 1.18+
- Git

###  Install and Run

```bash
# Clone the repo
git clone https://github.com/your-username/CountrySearch.git
cd CountrySearch

# Install dependencies
go mod tidy

# Run tests
go test ./... -v -cover

# Run the server
go run main.go
