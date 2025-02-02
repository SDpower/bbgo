package bbgo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c9s/bbgo/pkg/types"
)

func TestTradeCollector_ShouldNotCountDuplicatedTrade(t *testing.T) {
	symbol := "BTCUSDT"
	position := types.NewPosition(symbol, "BTC", "USDT")
	orderStore := NewOrderStore(symbol)
	collector := NewTradeCollector(symbol, position, orderStore)
	assert.NotNil(t, collector)

	matched := collector.ProcessTrade(types.Trade{
		ID:            1,
		OrderID:       399,
		Exchange:      types.ExchangeBinance,
		Price:         40000.0,
		Quantity:      1.0,
		QuoteQuantity: 40000.0,
		Symbol:        "BTCUSDT",
		Side:          types.SideTypeBuy,
		IsBuyer:       true,
	})
	assert.False(t, matched, "should be added to the trade store")
	assert.Equal(t, 1, len(collector.tradeStore.Trades()), "should have one trade in the trade store")

	orderStore.Add(types.Order{
		SubmitOrder:      types.SubmitOrder{
			Symbol:           "BTCUSDT",
			Side:             types.SideTypeBuy,
			Type:             types.OrderTypeLimit,
			Quantity:         1.0,
			Price:            40000.0,
		},
		Exchange:         types.ExchangeBinance,
		OrderID:          399,
		Status:           types.OrderStatusFilled,
		ExecutedQuantity: 1.0,
		IsWorking:        false,
	})

	matched = collector.Process()
	assert.True(t, matched)
	assert.Equal(t, 0, len(collector.tradeStore.Trades()), "the found trade should be removed from the trade store")

	matched = collector.ProcessTrade(types.Trade{
		ID:            1,
		OrderID:       399,
		Exchange:      types.ExchangeBinance,
		Price:         40000.0,
		Quantity:      1.0,
		QuoteQuantity: 40000.0,
		Symbol:        "BTCUSDT",
		Side:          types.SideTypeBuy,
		IsBuyer:       true,
	})
	assert.False(t, matched, "the same trade should not match")
	assert.Equal(t, 0, len(collector.tradeStore.Trades()), "the same trade should not be added to the trade store")
}
