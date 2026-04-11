package valuetypes

import "kumarvv.com/mockdata/utils"

const (
	SQL             = "sql"
	String          = "string"
	Number          = "number"
	Boolean         = "boolean"
	Date            = "date"
	DateTime        = "datetime"
	UUID            = "uuid"
	RandomString    = "random_string"
	RandomTitle     = "random_title"
	RandomGender    = "random_gender"
	RandomFirstName = "random_first_name"
	RandomLastName  = "random_last_name"
	RandomName      = "random_name"
	RandomEmail     = "random_email"
	RandomCurrency  = "random_currency"
	RandomAddress   = "random_address"
	RandomStreet    = "random_street"
	RandomCity      = "random_city"
	RandomState     = "random_state"
	RandomState2    = "random_state2"
	RandomCountry   = "random_country"
	RandomCountry2  = "random_country2"
	RandomCountry3  = "random_country3"
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
	RandomIn        = "random_in"
	RandomFrom      = "random_from"
)

func List() []string {
	return []string{
		SQL,
		String,
		Number,
		Boolean,
		Date,
		DateTime,
		UUID,
		RandomString,
		RandomTitle,
		RandomGender,
		RandomFirstName,
		RandomLastName,
		RandomName,
		RandomEmail,
		RandomCurrency,
		RandomAddress,
		RandomStreet,
		RandomCity,
		RandomState,
		RandomState2,
		RandomCountry,
		RandomCountry2,
		RandomCountry3,
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
		RandomIn,
		RandomFrom,
	}
}

func IsRequiredValueExpr(valueType string) bool {
	return utils.Includes([]string{
		SQL,
		String,
		Number,
		Boolean,
		Date,
		DateTime,
		RandomIn,
		RandomFrom,
	}, valueType)
}
