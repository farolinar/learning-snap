package snap

import (
	"fmt"
	"strconv"
)

// Service codes for the APIs covered in this lab.
const (
	ServiceAccessTokenB2B    = "73"
	ServiceBalanceInquiry    = "11"
	ServiceTransferIntrabank = "17"
)

// ResponseCode is the parsed form of a 7-digit SNAP response code.
type ResponseCode struct {
	HTTPStatus  int
	ServiceCode string
	CaseCode    string
}

// ComposeResponseCode builds a 7-digit code: HTTP status + 2-digit service + 2-digit case.
func ComposeResponseCode(httpStatus int, serviceCode, caseCode string) string {
	return fmt.Sprintf("%d%s%s", httpStatus, serviceCode, caseCode)
}

// ParseResponseCode splits a 7-digit code back into its parts.
func ParseResponseCode(code string) (ResponseCode, error) {
	if len(code) != 7 {
		return ResponseCode{}, fmt.Errorf("response code must be 7 digits, got %q", code)
	}
	httpStatus, err := strconv.Atoi(code[0:3])
	if err != nil {
		return ResponseCode{}, fmt.Errorf("invalid HTTP status in %q: %w", code, err)
	}
	return ResponseCode{
		HTTPStatus:  httpStatus,
		ServiceCode: code[3:5],
		CaseCode:    code[5:7],
	}, nil
}

// messages maps "httpStatus:caseCode" to the SNAP response message.
var messages = map[string]string{
	"200:00": "Successful",
	"202:00": "Request In Progress",
	"400:00": "Bad Request",
	"400:01": "Invalid Field Format",
	"400:02": "Invalid Mandatory Field",
	"401:00": "Unauthorized",
	"401:01": "Invalid Token (B2B)",
	"401:02": "Invalid Customer Token",
	"401:03": "Token Not Found (B2B)",
	"401:04": "Customer Token Not Found",
	"403:00": "Transaction Expired",
	"403:14": "Insufficient Funds",
	"403:18": "Inactive Card/Account/Customer",
	"404:00": "Invalid Transaction Status",
	"404:01": "Transaction Not Found",
	"404:11": "Invalid Card/Account/Customer",
	"404:18": "Inconsistent Request",
	"405:00": "Requested Function Is Not Supported",
	"409:00": "Conflict",
	"409:01": "Duplicate partnerReferenceNo",
	"429:00": "Too Many Requests",
	"500:00": "General Error",
	"500:01": "Internal Server Error",
	"504:00": "Timeout",
}

// MessageFor returns the SNAP response message for an HTTP status and case code.
func MessageFor(httpStatus int, caseCode string) string {
	if m, ok := messages[fmt.Sprintf("%d:%s", httpStatus, caseCode)]; ok {
		return m
	}
	return "Unknown"
}
