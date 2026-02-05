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

type Client struct {
	groupID    string
	httpClient *http.Client
}

func NewClient(groupID string) *Client {
	return &Client{
		groupID: groupID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type Participant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Currency     string        `json:"currency"`
	CurrencyCode string        `json:"currencyCode"`
	Participants []Participant `json:"participants"`
}

type Balance struct {
	Paid    int64 `json:"paid"`
	PaidFor int64 `json:"paidFor"`
	Total   int64 `json:"total"`
}

type Reimbursement struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

type BalanceResponse struct {
	Balances       map[string]Balance `json:"balances"`
	Reimbursements []Reimbursement    `json:"reimbursements"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PaidForEntry struct {
	Participant Participant `json:"participant"`
	Shares      int         `json:"shares"`
}

type Expense struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Amount      int64        `json:"amount"`
	ExpenseDate time.Time    `json:"expenseDate"`
	PaidBy      Participant  `json:"paidBy"`
	PaidFor     []PaidForEntry `json:"paidFor"`
	Category    *Category    `json:"category"`
}

type PaidForInput struct {
	ParticipantID string
	Shares        int
}

type ExpenseResult struct {
	ID string `json:"expenseId"`
}

type tRPCResult[T any] struct {
	Result struct {
		Data struct {
			JSON T `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

type tRPCBatchResponse[T any] []tRPCResult[T]

func (r tRPCBatchResponse[T]) First() (T, error) {
	var zero T
	if len(r) == 0 {
		return zero, fmt.Errorf("empty response")
	}
	return r[0].Result.Data.JSON, nil
}

func unmarshalTRPC[T any](data []byte) (T, error) {
	var response tRPCBatchResponse[T]
	var zero T
	if err := json.Unmarshal(data, &response); err != nil {
		return zero, err
	}
	return response.First()
}

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

func (c *Client) GetParticipants() ([]Participant, error) {
	group, err := c.GetGroup()
	if err != nil {
		return nil, err
	}
	return group.Participants, nil
}

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

func (c *Client) AddExpense(title, paidBy string, paidFor []PaidForInput, amount int64, category int, isReimbursement bool) (*ExpenseResult, error) {
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
					"isReimbursement":            isReimbursement,
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
