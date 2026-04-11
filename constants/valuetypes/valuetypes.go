package valuetypes

const (
	RandomString    = "random_string"
	RandomTitle     = "random_title"
	RandomGender    = "random_gender"
	RandomFirstName = "random_first_name"
	RandomLastName  = "random_last_name"
	RandomName      = "random_name"
	RandomEmail     = "random_email"
	RandomCountry   = "random_country"
	RandomCountry2  = "random_country2"
	RandomCountry3  = "random_country3"
	RandomCurrency  = "random_currency"
	RandomCity      = "random_city"
	RandomState     = "random_state"
	RandomState2    = "random_state2"
	RandomAddress   = "random_address"
	RandomNumber    = "random_number"
	RandomDecimal   = "random_decimal"
	RandomBoolean   = "random_boolean"
	RandomParagraph = "random_paragraph"
	RandomFormat    = "random_format"
	RandomDate      = "random_date"
	RandomDay       = "random_day"
	RandomMonth     = "random_month"
	RandomYear      = "random_year"
	RandomPhone     = "random_phone"
)

func List() []string {
	return []string{
		RandomString,
		RandomTitle,
		RandomGender,
		RandomFirstName,
		RandomLastName,
		RandomName,
		RandomEmail,
		RandomCountry,
		RandomCountry2,
		RandomCountry3,
		RandomCurrency,
		RandomCity,
		RandomState,
		RandomState2,
		RandomAddress,
		RandomNumber,
		RandomDecimal,
		RandomBoolean,
		RandomParagraph,
		RandomFormat,
		RandomDate,
		RandomDay,
		RandomMonth,
		RandomYear,
		RandomPhone,
	}
}
