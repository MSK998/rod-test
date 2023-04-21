package domain

import (
	"database/sql"
	"statement-service-poc/pkg/database"

	"github.com/MSK998/cube"
)

type StatementMerchantSummary struct {
	CardType      string
	NumberOfItems int
	CreditItems   int
	CreditCharges []byte
	ReturnItems   int
	ReturnCharges []byte
	NetCharges    []byte
	Sort          int
}

func GetSMSummary() ([]StatementMerchantSummary, error) {
	summ := make([]StatementMerchantSummary, 0)
	rows, err := database.Conn.Query("usp_GetStatementMerchantSummary", sql.Named("MerchantId", 226003), sql.Named("StatementPeriodId", 87))
	if err != nil {
		return nil, err
	}

	err = cube.ScanStruct(rows, &summ)
	if err != nil {
		return nil, err
	}

	// TODO MAP []byte to float
	return summ, nil
}
