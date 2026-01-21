package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://spliit.app/api/trpc"

// Client handles communication with the Spliit API
type Client struct {
	groupID    string
	httpClient *http.Client
}

// NewClient creates a new Spliit API client
func NewClient(groupID string) *Client {
	return &Client{
		groupID: groupID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Participant represents a group member
type Participant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Group represents a Spliit group
type Group struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Currency     string        `json:"currency"`
	CurrencyCode string        `json:"currencyCode"`
	Participants []Participant `json:"participants"`
}

// Balance represents a participant's balance
type Balance struct {
	Paid    int64 `json:"paid"`
	PaidFor int64 `json:"paidFor"`
	Total   int64 `json:"total"`
}

// Reimbursement represents a suggested payment
type Reimbursement struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

// BalanceResponse contains balances and reimbursements
type BalanceResponse struct {
	Balances       map[string]Balance `json:"balances"`
	Reimbursements []Reimbursement    `json:"reimbursements"`
}

// Category represents an expense category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PaidForEntry represents who an expense is split between
type PaidForEntry struct {
	Participant Participant `json:"participant"`
	Shares      int         `json:"shares"`
}

// Expense represents an expense entry
type Expense struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Amount      int64        `json:"amount"`
	ExpenseDate time.Time    `json:"expenseDate"`
	PaidBy      Participant  `json:"paidBy"`
	PaidFor     []PaidForEntry `json:"paidFor"`
	Category    *Category    `json:"category"`
}

// PaidForInput is used when creating expenses
type PaidForInput struct {
	ParticipantID string
	Shares        int
}

// ExpenseResult is returned after creating an expense
type ExpenseResult struct {
	ID string `json:"expenseId"`
}

// tRPCResult wraps a single tRPC response item
type tRPCResult[T any] struct {
	Result struct {
		Data struct {
			JSON T `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

// tRPCBatchResponse is an array of results (batch=1)
type tRPCBatchResponse[T any] []tRPCResult[T]

// First returns the first result's payload or error if empty
func (r tRPCBatchResponse[T]) First() (T, error) {
	var zero T
	if len(r) == 0 {
		return zero, fmt.Errorf("empty response")
	}
	return r[0].Result.Data.JSON, nil
}

// unmarshalTRPC parses a tRPC batch response and returns the first result
func unmarshalTRPC[T any](data []byte) (T, error) {
	var response tRPCBatchResponse[T]
	var zero T
	if err := json.Unmarshal(data, &response); err != nil {
		return zero, err
	}
	return response.First()
}

// Response payload types for tRPC endpoints
type groupResponse struct {
	Group Group `json:"group"`
}

type expensesResponse struct {
	Expenses []Expense `json:"expenses"`
}

type createGroupResponse struct {
	GroupID string `json:"groupId"`
}

func (c *Client) get(path string, params map[string]interface{}) ([]byte, error) {
	input, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("batch", "1")
	q.Set("input", string(input))
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *Client) post(path string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("batch", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetGroup retrieves group details
func (c *Client) GetGroup() (*Group, error) {
	params := map[string]interface{}{
		"0": map[string]interface{}{"json": map[string]string{"groupId": c.groupID}},
		"1": map[string]interface{}{"json": map[string]string{"groupId": c.groupID}},
	}

	body, err := c.get("/groups.get,groups.getDetails", params)
	if err != nil {
		return nil, err
	}

	result, err := unmarshalTRPC[groupResponse](body)
	if err != nil {
		return nil, err
	}

	return &result.Group, nil
}

// GetParticipants returns all participants in the group
func (c *Client) GetParticipants() ([]Participant, error) {
	group, err := c.GetGroup()
	if err != nil {
		return nil, err
	}
	return group.Participants, nil
}

// GetParticipantID finds a participant ID by name
func (c *Client) GetParticipantID(name string) (string, error) {
	participants, err := c.GetParticipants()
	if err != nil {
		return "", err
	}

	for _, p := range participants {
		if p.Name == name {
			return p.ID, nil
		}
	}

	return "", nil
}

// GetBalances retrieves balances and suggested reimbursements
func (c *Client) GetBalances() (*BalanceResponse, error) {
	params := map[string]interface{}{
		"0": map[string]interface{}{"json": map[string]string{"groupId": c.groupID}},
	}

	body, err := c.get("/groups.balances.list", params)
	if err != nil {
		return nil, err
	}

	result, err := unmarshalTRPC[BalanceResponse](body)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExpenses retrieves expenses for the group with pagination
func (c *Client) GetExpenses(limit, offset int) ([]Expense, error) {
	params := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"groupId": c.groupID,
				"cursor":  offset,
				"limit":   limit,
			},
		},
	}

	body, err := c.get("/groups.expenses.list", params)
	if err != nil {
		return nil, err
	}

	result, err := unmarshalTRPC[expensesResponse](body)
	if err != nil {
		return nil, err
	}

	return result.Expenses, nil
}

// AddExpense creates a new expense
func (c *Client) AddExpense(title, paidBy string, paidFor []PaidForInput, amount int64, category int) (*ExpenseResult, error) {
	paidForJSON := make([]map[string]interface{}, len(paidFor))
	for i, pf := range paidFor {
		paidForJSON[i] = map[string]interface{}{
			"participant": pf.ParticipantID,
			"shares":      pf.Shares,
		}
	}

	body := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"groupId": c.groupID,
				"expenseFormValues": map[string]interface{}{
					"expenseId":                  nil,
					"title":                      title,
					"expenseDate":                time.Now().Format(time.RFC3339),
					"amount":                     amount,
					"category":                   category,
					"paidBy":                     paidBy,
					"paidFor":                    paidForJSON,
					"isReimbursement":            false,
					"saveDefaultSplittingOptions": false,
					"notes":                      "",
					"documents":                  []interface{}{},
				},
			},
		},
	}

	respBody, err := c.post("/groups.expenses.create", body)
	if err != nil {
		return nil, err
	}

	result, err := unmarshalTRPC[ExpenseResult](respBody)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteExpense removes an expense
func (c *Client) DeleteExpense(expenseID string) error {
	body := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"groupId":   c.groupID,
				"expenseId": expenseID,
			},
		},
	}

	_, err := c.post("/groups.expenses.delete", body)
	return err
}

// CreateGroup creates a new group
func (c *Client) CreateGroup(name string, participants []string, currency, currencyCode string) (string, error) {
	participantObjs := make([]map[string]string, len(participants))
	for i, p := range participants {
		participantObjs[i] = map[string]string{"name": p}
	}

	body := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"groupFormValues": map[string]interface{}{
					"name":         name,
					"information":  "",
					"currency":     currency,
					"currencyCode": currencyCode,
					"participants": participantObjs,
				},
			},
		},
	}

	respBody, err := c.post("/groups.create", body)
	if err != nil {
		return "", err
	}

	result, err := unmarshalTRPC[createGroupResponse](respBody)
	if err != nil {
		return "", err
	}

	return result.GroupID, nil
}
