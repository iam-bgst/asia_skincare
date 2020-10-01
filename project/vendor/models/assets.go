package models

import (
	"addon"
	"time"
)

func InitialAssets() {
	delivery_model.InitialDelivery()
	membership_model.InitMembership()
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
