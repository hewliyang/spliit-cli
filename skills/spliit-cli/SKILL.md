---
name: spliit
description: "CLI tool for managing shared expenses using Spliit API. Track household expenses, balances, and reimbursements. Use when Claude needs to record expenses, check who owes whom, or manage group finances."
license: Custom
---

# Spliit CLI Tool

Manage shared expenses using the Spliit API.

## Configuration

**⚠️ ALWAYS check `memory/spliit.md` first** before running any Spliit command. This file contains all managed groups with their IDs, members, and linked Telegram GCs.

The group ID is passed via the `-g` or `--group` flag. You can find the group ID in the URL: `https://spliit.app/groups/<GROUP_ID>`

**When creating new groups**: Always update `memory/spliit.md` with:
- Group ID and URL
- Members
- Linked Telegram GC (if any)
- Currency

## Quick Commands

### Create a New Group
```bash
spliit create-group "Group Name" "Alice,Bob,Charlie"
spliit create-group "Trip" "Alice,Bob" --currency "$" --currency-code USD
```
Returns the new group ID and URL.

### Check Balances
```bash
spliit -g <group-id> balances
```
Shows who owes whom and suggested reimbursements.

### List Expenses
```bash
spliit -g <group-id> expenses
spliit -g <group-id> expenses --limit 10
spliit -g <group-id> expenses --limit 10 --page 2
spliit -g <group-id> expenses --from 2025-12-01
spliit -g <group-id> expenses --from 2025-12-01 --to 2025-12-31
spliit -g <group-id> expenses --show-ids   # Show IDs for deletion
```

### Add Expense
```bash
spliit -g <group-id> add-expense "Groceries" "Alice" 5000
spliit -g <group-id> add-expense "Movie tickets" "Bob" 3500 --category 1
```
**Note**: Amount is in cents (5000 = $50.00)

### Add Reimbursement (Settle a Debt)
```bash
spliit -g <group-id> add-expense "Settle up" "Bob" 2500 --reimbursement --to "Alice"
```
Use the `--reimbursement` flag with `--to` to record a debt settlement. This is used when someone pays back money they owe.

**Example workflow:**
1. Check balances: `spliit -g <group-id> balances`
   - Output shows: "Bob → Alice: $25.00"
2. Bob pays Alice $25 in cash/transfer
3. Record it: `spliit -g <group-id> add-expense "Settle up" "Bob" 2500 --reimbursement --to "Alice"`
4. Balances are now settled

**Important**: 
- The **payer** (positional arg) is the person paying back the debt
- The **--to** flag specifies who receives the money
- The `--to` flag is required when using `--reimbursement`

### Delete Expense
```bash
spliit -g <group-id> delete-expense <expense-id>
```
Use `expenses --show-ids` to get expense IDs first.

### List Participants
```bash
spliit -g <group-id> participants
```

### Get Group Details
```bash
spliit -g <group-id> group
```

## Important Notes

- **Amounts are in cents**: $50.00 = 5000 cents
- **Default split**: All expenses split equally among all participants
- **Categories**: 0=General, 1=Entertainment, 2=Food, 3=Transport
- **Reimbursements**: Use `--reimbursement` flag to record debt settlements

## Response Guidelines

When performing actions that produce URLs, always include them in your response:
- **Create group**: `https://spliit.app/groups/<GROUP_ID>`
- **Add expense**: `https://spliit.app/groups/<GROUP_ID>/expenses/<EXPENSE_ID>/edit` (direct link to the expense)
- **Balances**: `https://spliit.app/groups/<GROUP_ID>/balances`

This lets the user quickly verify or share the result.
