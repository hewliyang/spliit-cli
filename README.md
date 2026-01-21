# Spliit CLI

A command-line interface for managing shared expenses using the [Spliit](https://spliit.app) API.

## Installation

### Quick install (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/hewliyang/spliit-cli/main/install.sh | bash
```

### Download binary

Download the latest binary for your platform from [GitHub Releases](https://github.com/hewliyang/spliit-cli/releases/latest).

### Go install

```bash
go install github.com/hewliyang/spliit-cli@latest
```

### Build from source

```bash
git clone https://github.com/hewliyang/spliit-cli.git
cd spliit-cli
go build -o spliit
```

## Usage

All commands (except `create-group`) require the `--group` / `-g` flag.

You can find the group ID in the URL: `https://spliit.app/groups/<GROUP_ID>`

### Create Group

```bash
spliit create-group "Trip to Japan" "Alice,Bob,Charlie"

# With custom currency
spliit create-group "NYC Trip" "Alice,Bob" --currency "$" --currency-code USD
```

### Get Group Details

```bash
spliit -g <group-id> group
```

### List Participants

```bash
spliit -g <group-id> participants
```

### Check Balances

```bash
spliit -g <group-id> balances
```

### List Expenses

```bash
# List recent expenses (default: 50)
spliit -g <group-id> expenses

# Filter by date range
spliit -g <group-id> expenses --from 2025-11-01 --to 2025-11-30

# Limit results and show IDs
spliit -g <group-id> expenses --limit 10 --show-ids
```

### Add Expense

```bash
# Amount is in cents (5000 = $50.00)
spliit -g <group-id> add-expense "Groceries" "Alice" 5000

# With category
spliit -g <group-id> add-expense "Movie" "Bob" 3500 --category 1
```

Categories:
- 0 = General (default)
- 1 = Entertainment
- 2 = Food
- 3 = Transport

### Delete Expense

```bash
spliit -g <group-id> delete-expense <expense-id>
```

Use `expenses --show-ids` to find expense IDs.

## Example Workflow

```bash
# Create a group
spliit create-group "Weekend Trip" "Alice,Bob,Charlie"
# Returns: Group ID: abc123

# Add expenses
spliit -g abc123 add-expense "Hotel" "Alice" 30000
spliit -g abc123 add-expense "Dinner" "Bob" 8500

# Check who owes whom
spliit -g abc123 balances
```

## License

MIT
