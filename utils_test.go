package gengormstruct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFmtFieldName(t *testing.T) {
	testCases := map[string]string{
		"hello_world":   "HelloWorld",
		"order_uid":     "OrderUid",
		"Order_Uid":     "OrderUid",
		"_order_uid":    "OrderUid",
		"_order3_uid":   "Order3Uid",
		"_order_3_uid":  "Order3Uid",
		"_order_3_uid_": "Order3Uid",
		"3order_3_uid_": "_order3Uid",
	}

	for key, value := range testCases {
		assert.Equal(t, value, CamelCaseString(key))
	}
}

func TestFilterColumnTypeSize(t *testing.T) {
	testCases := map[string]string{
		"bigint(21) unsigned": "bigint unsigned",
		"varchar(1024)":       "varchar",
		"datetime":            "datetime",
	}
	for key, value := range testCases {
		assert.Equal(t, value, FilterColumnTypeSize(key))
	}
}
