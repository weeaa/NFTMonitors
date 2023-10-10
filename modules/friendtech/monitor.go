package friendtech

import (
	"context"
	"encoding/hex"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/utils"
	"math/big"
	"net/url"
	"strings"
)

func (s *Settings) MonitorFriendTechLogs() {
	ch := make(chan types.Log)
	sub, err := utils.CreateSubscription(s.WSSClient, []common.Address{common.HexToAddress(FRIEND_TECH_CONTRACT_V1)}, ch)
	if err != nil {
		logger.LogFatal(moduleName, err)
	}

	for {
		select {
		case <-s.Context.Done():
			logger.LogShutDown(moduleName)
		case err = <-sub.Err():
			logger.LogError(moduleName, fmt.Errorf("subscription stopped: %w", err))
			return
		case txns := <-ch:
			var tx *types.Transaction
			var sender common.Address
			var balance *big.Int
			_ = balance
			tx, _, err = s.WSSClient.TransactionByHash(context.Background(), txns.TxHash)
			if err != nil {
				logger.LogError(moduleName, err)
				continue
			}

			sender, err = utils.GetSender(tx)
			if err != nil {
				logger.LogError(moduleName, err)
				continue
			}

			// get the sender's balance
			balance, err = utils.GetEthWalletBalance(s.WSSClient, sender)
			if err != nil {
				logger.LogError(moduleName, err)
				continue
			}
		}
	}
}

// getLatestUser will be used to determine what is the latest ID active.
func (s *Settings) getLatestUser(startUserID int) bool {

	// would go faster to check by chunks rather than 1 by 1
	steps := []int{5000, 3000, 1000, 500, 100}

	for _, step := range steps {
		userID := startUserID

		for {
			req := &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Scheme: "https", Host: prodBaseApi, Path: "/users/by-id/" + fmt.Sprint(userID)},
				Host:   prodBaseApi,
				Header: http.Header{},
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				break
			}

			if resp.StatusCode == http.StatusOK {
				userID += step
			} else if resp.StatusCode == http.StatusNotFound {
				userID -= step
				break
			} else {
				fmt.Println("Unexpected status code:", resp.StatusCode)
				break
			}

			if err = resp.Body.Close(); err != nil {
				fmt.Println("Error closing response body:", err)
				break
			}
		}

		// After completing the inner loop, we have the latest valid user ID for this step
		// If the step is 100, we can refine the search by checking smaller increments
		if step == 100 {
			for i := 1; i <= 100; i++ {
				req := &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Scheme: "https", Host: prodBaseApi, Path: "/users/by-id/" + fmt.Sprint(userID+i)},
					Host:   prodBaseApi,
					Header: http.Header{},
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					// Handle error
					fmt.Println("Error:", err)
					break
				}

				if resp.StatusCode == http.StatusOK {
					userID += i
				} else {
					break
				}

				if err := resp.Body.Close(); err != nil {
					// Handle error when closing response body
					fmt.Println("Error closing response body:", err)
					break
				}
			}
		}
	}

	// Return the latest valid user ID found
	return false
}

func (s *Settings) dispatchLog(txns types.Log) {
	var tx *types.Transaction
	var sender common.Address
	var balance *big.Int
	var err error

	_ = balance

	tx, _, err = s.WSSClient.TransactionByHash(context.Background(), txns.TxHash)
	if err != nil {
		logger.LogError(moduleName, err)
	}

	sender, err = utils.GetSender(tx)
	if err != nil {
		logger.LogError(moduleName, err)
	}

	balance, err = utils.GetEthWalletBalance(s.WSSClient, sender)
	if err != nil {
		logger.LogError(moduleName, err)
	}

	if strings.Contains(tx.Hash().Hex(), sellMethod) {
		if err = s.handleSell(tx, sender.String()); err != nil {
			logger.LogError(moduleName, err)
		}
	} else if strings.Contains(tx.Hash().Hex(), buyMethod) {
		s.handleBuy(tx, sender.String())
	} else {

	}
}

func (s *Settings) monitorBalance() {
	go func() {

	}()
}

func (s *Settings) handleSell(tx *types.Transaction, sender string) error {
	/*
		if isSelf(hex.EncodeToString(tx.Data()), sender) {
			user, err := s.retrieveUserFromDB(sender)
			if err != nil {
				return err
			}
			imp, err := assertImportance() //pass followers count
			// do a go func that adds to the database the txn data
			return s.Discord.SendNotification(buildWebhook(), "")
		} else {

		}
	*/
	return nil
}

func (s *Settings) handleBuy(tx *types.Transaction, sender string) {
	if isSelf(hex.EncodeToString(tx.Data()), sender) {

	} else {

	}
}

func isSelf(txData string, sender string) bool {
	return strings.Contains(strings.ToLower(txData), strings.ToLower(sender[2:]))
}

// TODO DOES NOT WORK – cannot confirm if it's really a new user, would need to fetch tx history
// isNewUser validates whether a user has just signed up or not.
func isNewUser(txData, sender string) bool {
	return strings.Contains(strings.ToLower(txData), strings.ToLower(sender[2:])) && txData[len(txData)-1:] == "1"
}
