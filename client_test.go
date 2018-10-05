package iex

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockHTTPClient struct {
	body    string
	headers map[string]string
	code    int
	err     error
}

func (c *mockHTTPClient) Get(url string) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteString(c.body)

	for key, value := range c.headers {
		w.Header().Add(key, value)
	}

	w.WriteHeader(c.code)

	resp := w.Result()
	return resp, c.err
}

func setupTestClient() *Client {
	return NewClient(&http.Client{
		Timeout: 5 * time.Second,
	})
}

func TestTOPS_AllSymbols(t *testing.T) {
	// TODO: Add expected field to struct and use it to verify results
	var testCases = []struct {
		symbols []string
		code    int
		body    string
		err     error
		headers map[string]string
	}{
		{symbols: []string{"SNAP", "FB"}, code: 200, body: `[{"symbol":"SNAP","sector":"softwareservices","securityType":"commonstock","bidPrice":0,"bidSize":0,"askPrice":0,"askSize":0,"lastUpdated":1537215438021,"lastSalePrice":9.165,"lastSaleSize":123,"lastSaleTime":1537214395927,"volume":525079,"marketPercent":0.0238},{"symbol":"FB","sector":"softwareservices","securityType":"commonstock","bidPrice":0,"bidSize":0,"askPrice":0,"askSize":0,"lastUpdated":1537216916977,"lastSalePrice":160.6,"lastSaleSize":100,"lastSaleTime":1537214399372,"volume":991898,"marketPercent":0.04741}]`, err: nil, headers: map[string]string{"Content-Type": "application/json"}},
		{symbols: []string{"AIG+"}, code: 200, body: `[{"symbol":"AIG+","sector":"n/a","securityType":"warrant","bidPrice":0,"bidSize":0,"askPrice":0,"askSize":0,"lastUpdated":1537214400001,"lastSalePrice":0,"lastSaleSize":0,"lastSaleTime":0,"volume":0,"marketPercent":0}]`, err: nil, headers: map[string]string{"Content-Type": "application/json"}},
	}

	for _, tt := range testCases {
		httpc := mockHTTPClient{body: tt.body, code: tt.code, err: tt.err, headers: tt.headers}
		c := NewClient(&httpc)

		result, err := c.GetTOPS(tt.symbols)

		if err != nil {
			t.Fatal(err)
		}

		if len(result) != len(tt.symbols) {
			t.Fatalf("Received %v results, expected %v", len(result), len(tt.symbols))
		}
	}
}

func TestLast(t *testing.T) {
	c := setupTestClient()
	symbols := []string{"SPY", "AAPL"}
	result, err := c.GetLast(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Received %v results, expected %v", len(result), len(symbols))
	}
}

func TestHIST_OneDate(t *testing.T) {
	c := setupTestClient()
	testDate := time.Date(2017, time.June, 6, 0, 0, 0, 0, time.UTC)
	result, err := c.GetHIST(testDate)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatalf("Received zero results")
	}
}

func TestHIST_AllDates(t *testing.T) {
	c := setupTestClient()
	result, err := c.GetAllAvailableHIST()
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatalf("Received zero results")
	}
}

func TestDEEP(t *testing.T) {
	c := setupTestClient()
	result, err := c.GetDEEP("SPY")
	if err != nil {
		t.Fatal(err)
	}

	if result.Symbol != "SPY" {
		t.Fatalf("Expected symbol = %v, got %v", "SPY", result.Symbol)
	}
}

func TestBook(t *testing.T) {
	body := `{
		"YELP": {
			"bids": [
				{
					"price": 63.09,
					"size": 300,
					"timestamp": 1494538496261
				}
			],
			"asks": [
				{
					"price": 63.92,
					"size": 300,
					"timestamp": 1494538381896
				},
				{
					"price": 63.97,
					"size": 300,
					"timestamp": 1494538381885
				}
			]
		}
	}`
	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"SPY"}
	result, err := c.GetBook(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Log(result)
		t.Fatalf("Received %v results, expected %v", len(result), len(symbols))
	}
}

func TestGetTrades(t *testing.T) {
	body := `{
	"AAPL": [],
	"FB": []
}`
	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}
	last := 1

	result, err := c.GetTrades(symbols, last)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestGetSystemEvents(t *testing.T) {
	body := `{
	"AAPL": {
		"systemEvent": "R",
		"timestamp": 1494627280251
	},
	"FB": {
		"systemEvent": "R",
		"timestamp": 1494627280251
	}
}`

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}

	result, err := c.GetSystemEvents(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestGetTradingStatus(t *testing.T) {
	body := `{
	"AAPL": {
		"status": "T",
    "reason": "NA",
    "timestamp": 1494588017687
	},
	"FB": {
		"status": "T",
    "reason": "NA",
    "timestamp": 1494588017687
	}
}`

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}

	result, err := c.GetTradingStatus(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestGetOperationalHaltStatus(t *testing.T) {
	body := `{
	"AAPL": {
		"isHalted": false,
    "timestamp": 1494588017687
	},
	"FB": {
		"isHalted": false,
    "timestamp": 1494588017687
	}
}`

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}

	result, err := c.GetOperationalHaltStatus(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestGetShortSaleRestriction(t *testing.T) {
	body := `{
	"AAPL": {
    "isSSR": true,
    "detail": "N",
    "timestamp": 1494588094067
	},
	"FB": {
    "isSSR": true,
    "detail": "N",
    "timestamp": 1494588094067
	}
}`

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}

	result, err := c.GetShortSaleRestriction(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestGetSecurityEvents(t *testing.T) {
	body := `{
	"AAPL": {
    "securityEvent": "MarketOpen",
    "timestamp": 1494595800005
	},
	"FB": {
    "securityEvent": "MarketOpen",
    "timestamp": 1494595800005
	}
}`

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}

	result, err := c.GetSecurityEvents(symbols)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestTradeBreaks(t *testing.T) {
	body := `{
	"AAPL": [
		{
      "price": 156.1,
      "size": 100,
      "tradeId": 517341294,
      "isISO": false,
      "isOddLot": false,
      "isOutsideRegularHours": false,
      "isSinglePriceCross": false,
      "isTradeThroughExempt": false,
      "timestamp": 1494619192003
		}
	],
	"FB": [
 		{
      "price": 156.1,
      "size": 100,
      "tradeId": 517341294,
      "isISO": false,
      "isOddLot": false,
      "isOutsideRegularHours": false,
      "isSinglePriceCross": false,
      "isTradeThroughExempt": false,
      "timestamp": 1494619192003
		}
	]
}`

	httpc := mockHTTPClient{body: body, code: 200}
	c := NewClient(&httpc)

	symbols := []string{"AAPL", "FB"}
	last := 1

	result, err := c.GetTradeBreaks(symbols, last)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(symbols) {
		t.Fatalf("Number of symbols returned %d, not equal to requested %d",
			len(result), len(symbols))
	}
}

func TestSymbols(t *testing.T) {
	c := setupTestClient()
	symbols, err := c.GetSymbols()
	if err != nil {
		t.Fatal(err)
	}

	if len(symbols) == 0 {
		t.Fatal("Received zero symbols")
	}

	symbol := symbols[0]
	if symbol.Symbol == "" || symbol.Name == "" || symbol.Date == "" {
		t.Fatal("Failed to decode symbol correctly")
	}
}

func TestMarkets(t *testing.T) {
	c := setupTestClient()
	markets, err := c.GetMarkets()
	if err != nil {
		t.Fatal(err)
	}

	if len(markets) == 0 {
		t.Fatal("Received zero markets")
	}
}

func TestGetHistoricalDaily(t *testing.T) {
	c := setupTestClient()
	stats, err := c.GetHistoricalDaily(&HistoricalDailyRequest{Last: 5})
	if err != nil {
		t.Fatal(err)
	}

	if len(stats) != 5 {
		t.Fatalf("Received %d historical daily stats, expected %d", len(stats), 5)
	}
}
