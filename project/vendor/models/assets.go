package models

import (
	"addon"
	"sync"
	"time"
)

func InitialAssets() {
	delivery_model.InitialDelivery()
	membership_model.InitMembership()
	payment_model.InitialPayment()
	courier_model.InitialCourier()
}

func ServicePoint() {
	go func() {
		for {
			account := account_model.All()
			for _, a := range account {
				exp := a.Point.Exp.AddDate(2, 0, 0)
				// fmt.Println("A", a.RegiteredAt)
				// fmt.Println("B", registeredAt)
				// fmt.Println("same", addon.DateSameOrNot(registeredAt, a.RegiteredAt))
				if addon.DateSameOrNot(exp, time.Now()) {
					account_model.UpdatePoint(a.Id, a.Point.Value-(a.Point.Value*2))
					account_model.UpdateExpPoint(a.Id, exp)
				} else {
					// fmt.Println("else")
					continue
				}
			}
		}
	}()
}

func ServiceTransaction() {
	go func() {
		wg := &sync.WaitGroup{}
		for {
			wg.Add(1)
			data := transaction_model.All()
			for _, d := range data {
				sub := time.Now().Sub(d.Date).Hours()
				if sub > 24.0 && d.Evidence == (Evidence{}) && d.Status_code == 0 {
					// fmt.Println(d.Id)
					transaction_model.UpdateStatus(d.Id, CENCELED)
				} else {
					continue
				}
			}
			wg.Done()
			wg.Wait()
		}
	}()
}
