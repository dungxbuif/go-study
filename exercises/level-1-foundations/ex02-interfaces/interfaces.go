package main

import (
	"errors"
	"fmt"
	"time"
)

var StripeAmountExceededError = errors.New("stripe: amount exceeds limit")
var MomoCurrencyNotSupportError = errors.New("momo: only supports VND")

type PaymentGateway interface {
	Charge(amount int64, currency string) (string, error)
	Refund(txID string) error
}

type StripeGateway struct{}

func (s StripeGateway) Charge(amount int64, currency string) (string, error) {
	if amount > 10000 {
		return "", StripeAmountExceededError
	}
	txID := fmt.Sprintf("stripe_%d", time.Now().UnixMilli())
	fmt.Printf("Stripe: Charged %d %s -> tx: %s\n", amount, currency, txID)
	return txID, nil
}

func (s StripeGateway) Refund(txID string) error {
	fmt.Printf("Stripe: Refunded tx %s\n", txID)
	return nil
}

type MomoGateway struct{}

func (s MomoGateway) Charge(amount int64, currency string) (string, error) {
	if currency != "VND" {
		return "", MomoCurrencyNotSupportError
	}
	txID := fmt.Sprintf("momo_%d", time.Now().UnixMilli())
	fmt.Printf("Momo: Charged %d %s -> tx: %s\n", amount, currency, txID)
	return txID, nil
}

func (s MomoGateway) Refund(txID string) error {
	if txID == "" {
		return errors.New("momo: txID is required for refund")
	}
	fmt.Printf("Momo: Refunded tx %s\n", txID)
	return nil
}

func processOrder(gw PaymentGateway, amount int64,) {
	txID, err := gw.Charge(amount, "VND")
	if err != nil {
		fmt.Printf("order failed: %v\n", err)
		return
	}
	fmt.Printf("Order processed successfully: %s\n", txID)
}

func main() {
	stripe := StripeGateway{}
	momo := MomoGateway{}
	processOrder(stripe, 5000)
	processOrder(stripe, 15000)
	processOrder(momo, 50000)
}
