package debankl2

import (
	"debank_checker_v3/global"
	"debank_checker_v3/util"
	"fmt"
	"github.com/sirupsen/logrus"
)

type L2ParserDebank struct{}

func (d L2ParserDebank) Parse(accountData string) {
	accountData = util.RemoveHexPrefix(accountData)

	address, accountType, privKey, err := util.GetAccountData(accountData)
	if err != nil {
		logrus.Printf("[%d/%d] | %s | Failed to parse account data: %v",
			global.CurrentProgress, global.TargetProgress, accountData, err)
		return
	}

	if accountType != 1 && accountType != 2 {
		logrus.Printf("[%d/%d] | %s | Not a private key or mnemonic",
			global.CurrentProgress, global.TargetProgress, accountData)
		return
	}

	message := getSignL2(address)
	signature := signMessage(message, privKey)
	captcha := solveCaptcha(address)
	sessionID := doAuth(address, captcha, signature)
	balance := getL2Balance(address, sessionID)

	logrus.Printf("[%d/%d] | %s | L2 Balance: $%.2f",
		global.CurrentProgress, global.TargetProgress, address, balance)

	if balance > 0 {
		util.AppendToFile("./results/debank_l2_balances.txt",
			fmt.Sprintf("%s | %s | $%.2f\n", accountData, address, balance))
	}
}
