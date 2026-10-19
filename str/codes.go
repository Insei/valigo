package str

import "strings"

// countryCodes are the officially assigned codes of ISO 3166-1 alpha-2.
var countryCodes = codeSet(`
AD AE AF AG AI AL AM AO AQ AR AS AT AU AW AX AZ BA BB BD BE BF BG BH BI BJ BL BM BN BO BQ BR BS BT BV BW BY BZ
CA CC CD CF CG CH CI CK CL CM CN CO CR CU CV CW CX CY CZ DE DJ DK DM DO DZ EC EE EG EH ER ES ET FI FJ FK FM FO FR
GA GB GD GE GF GG GH GI GL GM GN GP GQ GR GS GT GU GW GY HK HM HN HR HT HU ID IE IL IM IN IO IQ IR IS IT JE JM JO JP
KE KG KH KI KM KN KP KR KW KY KZ LA LB LC LI LK LR LS LT LU LV LY MA MC MD ME MF MG MH MK ML MM MN MO MP MQ MR MS MT
MU MV MW MX MY MZ NA NC NE NF NG NI NL NO NP NR NU NZ OM PA PE PF PG PH PK PL PM PN PR PS PT PW PY QA RE RO RS RU RW
SA SB SC SD SE SG SH SI SJ SK SL SM SN SO SR SS ST SV SX SY SZ TC TD TF TG TH TJ TK TL TM TN TO TR TT TV TW TZ UA UG
UM US UY UZ VA VC VE VG VI VN VU WF WS YE YT ZA ZM ZW`)

// currencyCodes are the active codes of ISO 4217 in 2026, with the codes of funds and precious metals.
var currencyCodes = codeSet(`
AED AFN ALL AMD AOA ARS AUD AWG AZN BAM BBD BDT BHD BIF BMD BND BOB BOV BRL BSD BTN BWP BYN BZD CAD CDF CHE CHF CHW
CLF CLP CNY COP COU CRC CUP CVE CZK DJF DKK DOP DZD EGP ERN ETB EUR FJD FKP GBP GEL GHS GIP GMD GNF GTQ GYD HKD HNL
HTG HUF IDR ILS INR IQD IRR ISK JMD JOD JPY KES KGS KHR KMF KPW KRW KWD KYD KZT LAK LBP LKR LRD LSL LYD MAD MDL MGA
MKD MMK MNT MOP MRU MUR MVR MWK MXN MXV MYR MZN NAD NGN NIO NOK NPR NZD OMR PAB PEN PGK PHP PKR PLN PYG QAR RON RSD
RUB RWF SAR SBD SCR SDG SEK SGD SHP SLE SOS SRD SSP STN SVC SYP SZL THB TJS TMT TND TOP TRY TTD TWD TZS UAH UGX USD
USN UYI UYU UYW UZS VED VES VND VUV WST XAF XAG XAU XBA XBB XBC XBD XCD XCG XDR XOF XPD XPF XPT XSU XTS XUA XXX YER
ZAR ZMW ZWG`)

func codeSet(codes string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, code := range strings.Fields(codes) {
		set[code] = struct{}{}
	}
	return set
}

// isCountryCode reports whether s is an upper case code of ISO 3166-1 alpha-2: US, DE.
func isCountryCode(s string) bool {
	_, ok := countryCodes[s]
	return ok
}

// isCurrencyCode reports whether s is an upper case code of ISO 4217: USD, EUR.
func isCurrencyCode(s string) bool {
	_, ok := currencyCodes[s]
	return ok
}

// isLuhn reports whether s is at least 2 digits with a valid checksum of the Luhn algorithm.
func isLuhn(s string) bool {
	return len(s) >= 2 && luhn(s, false)
}

// isCreditCard reports whether s is a number of 12-19 digits with a valid Luhn checksum,
// the digits can be grouped by spaces or hyphens: 4111 1111 1111 1111.
func isCreditCard(s string) bool {
	digits := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case isASCIIDigit(c):
			digits++
		case (c == ' ' || c == '-') && i > 0 && i < len(s)-1 && isASCIIDigit(s[i-1]):
		default:
			return false
		}
	}
	return digits >= 12 && digits <= 19 && luhn(s, true)
}

// luhn checks the checksum of the digits of s by the Luhn algorithm, separators are skipped if allowed.
func luhn(s string, separators bool) bool {
	sum, double := 0, false
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if !isASCIIDigit(c) {
			if separators {
				continue
			}
			return false
		}
		d := int(c - '0')
		if double {
			if d *= 2; d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}
