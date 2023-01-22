package internal

import "strconv"

type AccountType uint8

const (
	SavingAccount AccountType = iota + 1
	GiroAccount
)

var (
	AccountTypeValue = map[AccountType]string{
		SavingAccount: "S", // Saving
		GiroAccount:   "G", // Giro
	}
)

func (x AccountType) String() string {
	return AccountTypeValue[x]
}

type AccountStatus uint8

const (
	StatusOpen AccountStatus = iota + 1
	StatusClose
)

var (
	AccountStatusValue = map[AccountStatus]string{
		StatusOpen:  "1",
		StatusClose: "2",
	}
)

func (x AccountStatus) String() string {
	return AccountStatusValue[x]
}

type AccountIsBlocked uint8

const (
	NotBlocked AccountIsBlocked = iota
	Blocked
)

var (
	AccountIsBlockedValue = map[AccountIsBlocked]string{
		NotBlocked: "0",
		Blocked:    "1",
	}
)

func (x AccountIsBlocked) String() string {
	return AccountIsBlockedValue[x]
}

func BuildNewCif(lastCif string) string {
	lastCifNum, _ := strconv.Atoi(lastCif)
	var lastCifNumDigits int
	for n := lastCifNum; n > 0; n-- {
		n /= 10
		lastCifNumDigits++
	}
	var newCif string
	for i := 10; i > lastCifNumDigits; i-- {
		newCif += "0"
	}
	lastCifNum++
	newCif += strconv.Itoa(lastCifNum)
	return newCif
}

const (
	AccountHead = "00100100"
	AccountTail = "300"
)

func BuildNewAccountNumber(lastAccount string) string {
	lastAccountId, _ := strconv.Atoi(lastAccount[8:12])
	newAccountId := lastAccountId + 1
	var newAccountDigits int
	for n := newAccountId; n > 0; n-- {
		n /= 10
		newAccountDigits++
	}
	var newAccount string
	for i := 4; i > newAccountDigits; i-- {
		newAccount += "0"
	}
	newAccount += strconv.Itoa(newAccountId)
	return AccountHead + newAccount + AccountTail
}
