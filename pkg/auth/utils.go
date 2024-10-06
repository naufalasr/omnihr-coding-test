package auth

import (
	"fmt"
	"strconv"
	"strings"
)

const separation = ":"

func GenerateTokenParam(userName string, companyID int64) string {
	return fmt.Sprintf("%s%s%d", userName, separation, companyID)
}

func ParseClaimUsername(token string) (userName string, companyID int64) {
	res := strings.Split(token, separation)
	userName = res[0]
	companyID, _ = strconv.ParseInt(res[1], 0, 64)
	return userName, companyID
}
