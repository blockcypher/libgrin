package client_test

import (
	"fmt"
	"testing"

	"github.com/blockcypher/libgrin/v5/client"
	"github.com/stretchr/testify/assert"
)

func TestWalletForeignAPI(t *testing.T) {
	// commenting this since this can't be done on CI for now
	url := "http://127.0.0.1:3415/v2/foreign"
	walletForeignAPI := client.NewWalletForeignAPI(url)
	// CheckVersion
	{
		versionInfo, err := walletForeignAPI.CheckVersion()
		fmt.Println(versionInfo)
		assert.NoError(t, err)
		assert.NotNil(t, versionInfo)
	}
}
