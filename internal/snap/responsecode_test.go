package snap

import "testing"

func TestComposeParseRoundTrip(t *testing.T) {
	cases := []struct {
		httpStatus int
		service    string
		caseCode   string
		want       string
	}{
		{200, ServiceAccessTokenB2B, "00", "2007300"},
		{200, ServiceBalanceInquiry, "00", "2001100"},
		{200, ServiceTransferIntrabank, "00", "2001700"},
		{400, ServiceBalanceInquiry, "02", "4001102"},
		{401, ServiceBalanceInquiry, "01", "4011101"},
		{403, ServiceBalanceInquiry, "14", "4031114"},
		{404, ServiceBalanceInquiry, "18", "4041118"},
		{409, ServiceTransferIntrabank, "00", "4091700"},
	}
	for _, c := range cases {
		got := ComposeResponseCode(c.httpStatus, c.service, c.caseCode)
		if got != c.want {
			t.Errorf("ComposeResponseCode(%d,%s,%s)=%q want %q", c.httpStatus, c.service, c.caseCode, got, c.want)
		}
		parsed, err := ParseResponseCode(got)
		if err != nil {
			t.Fatalf("ParseResponseCode(%q): %v", got, err)
		}
		if parsed.HTTPStatus != c.httpStatus || parsed.ServiceCode != c.service || parsed.CaseCode != c.caseCode {
			t.Errorf("parse round trip failed for %q: %+v", got, parsed)
		}
	}
}

func TestParseResponseCodeRejectsBadLength(t *testing.T) {
	if _, err := ParseResponseCode("123"); err == nil {
		t.Fatal("expected error for short code")
	}
}

func TestMessageFor(t *testing.T) {
	if got := MessageFor(200, "00"); got != "Successful" {
		t.Errorf("200/00 = %q, want Successful", got)
	}
	if got := MessageFor(403, "14"); got != "Insufficient Funds" {
		t.Errorf("403/14 = %q, want Insufficient Funds", got)
	}
	if got := MessageFor(999, "99"); got != "Unknown" {
		t.Errorf("unknown = %q, want Unknown", got)
	}
}
