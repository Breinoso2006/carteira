package models

import (
	"encoding/json"
	"math"
	"time"

	"github.com/breinoso2006/carteira-api/internal/http"
)

type Stock struct {
	Ticker              string
	FundamentalistGrade float64
	Moment              float64
	FinalGrade          float64
}

type StockResponse struct {
	Symbol string  `json:"Symbol"`
	Price  float64 `json:"Price"`
	PE     float64 `json:"PE"`
	PBV    float64 `json:"PBV"`
	PSR    float64 `json:"PSR"`
	BVps   float64 `json:"BVps"`
	EPS    float64 `json:"EPS"`
	DY     float64 `json:"DY"`
	Source string  `json:"Source"`
}

func NewStock(ticker string, fundamentalistGrade float64) *Stock {
	return &Stock{
		Ticker:              ticker,
		FundamentalistGrade: fundamentalistGrade,
	}
}

func (s *Stock) SetFinalGrade() {
	if err := s.setMoment(); err != nil {
		return
	}

	alpha := 0.1

	momentNorm := s.Moment / 6.0

	s.FinalGrade = adjustScore(s.FundamentalistGrade * (1 + alpha*(momentNorm-0.5)))
}

func adjustScore(score float64) float64 {
	normalized := score / 100.0

	if normalized > 0.75 {
		normalized *= 1.08
	}
	if normalized < 0.65 {
		normalized *= 0.92
	}

	return math.Pow(normalized, 1.4)
}

func (s *Stock) setMoment() error {
	client := http.NewHTTPClient(5 * time.Second)
	body, err := client.Get("http://localhost:3001/" + s.Ticker)
	if err != nil {
		return err
	}

	var stockMomentData StockResponse
	if err := json.Unmarshal(body, &stockMomentData); err != nil {
		return err
	}

	s.calculateMoment(&stockMomentData)

	return nil
}

func (s *Stock) calculateMoment(stockData *StockResponse) {
	moment := 0

	moment += isPeGood(stockData.PE)
	moment += isPbvGood(stockData.PBV)
	moment += isPsrGood(stockData.PSR)
	moment += isDyGood(stockData.DY)
	moment += isGrahamGood(stockData.Price, stockData.EPS, stockData.BVps)
	s.Moment = float64(moment)
}

func isPeGood(pe float64) int {
	if pe > 0 && pe <= 8 {
		return 1
	}
	return 0
}

func isPbvGood(pbv float64) int {
	if pbv > 0 && pbv <= 2 {
		return 1
	}
	return 0
}

func isPsrGood(psr float64) int {
	if psr > 0 && psr < 2 {
		return 1
	}
	return 0
}

func isDyGood(dy float64) int {
	if dy >= 4 {
		return 1
	}
	return 0
}

func isGrahamGood(price, eps, bvps float64) int {
	moment := 0
	if eps > 0 && bvps > 0 {
		grahamValue := math.Sqrt(22.5 * eps * bvps)
		if math.IsNaN(grahamValue) {
			return 0
		}
		if grahamValue > 0 && price < grahamValue {
			moment += 1
			if grahamValue < price*0.7 {
				moment += 1
			}
		}
	}
	return moment
}
