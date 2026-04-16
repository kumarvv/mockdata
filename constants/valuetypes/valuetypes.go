package valuetypes

import "kumarvv.com/mockdata/utils"

const (
	SQL             = "sql"
	String          = "string"
	Integer         = "integer"
	Float           = "float"
	Boolean         = "boolean"
	Date            = "date"
	DateTime        = "datetime"
	Serial          = "serial"
	UUID            = "uuid"
	RandomString    = "random_string"
	RandomTitle     = "random_title"
	RandomGender    = "random_gender"
	RandomFirstName = "random_first_name"
	RandomLastName  = "random_last_name"
	RandomFullName  = "random_full_name"
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
	RandomInString  = "random_in_string"
	RandomInInteger = "random_in_integer"
	RandomInFloat   = "random_in_float"
	RandomRange     = "random_range"
	RandomFromSQL   = "random_from_sql"
)

func List() []string {
	return []string{
		SQL,
		String,
		Integer,
		Float,
		Boolean,
		Date,
		DateTime,
		Serial,
		UUID,
		RandomString,
		RandomTitle,
		RandomGender,
		RandomFirstName,
		RandomLastName,
		RandomFullName,
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
		RandomInString,
		RandomInInteger,
		RandomInFloat,
		RandomRange,
		RandomFromSQL,
	}
}

func IsRequiredValueExpr(valueType string) bool {
	return utils.Includes([]string{
		SQL,
		String,
		Integer,
		Float,
		Boolean,
		Date,
		DateTime,
		RandomInString,
		RandomInInteger,
		RandomInFloat,
		RandomFromSQL,
	}, valueType)
}

func IsDbRequired(valueType string) bool {
	return utils.Includes([]string{
		SQL,
		RandomFromSQL,
	}, valueType)
}

func IsString(valueType string) bool {
	return utils.Includes([]string{
		String,
		RandomString,
		RandomTitle,
		RandomGender,
		RandomFirstName,
		RandomLastName,
		RandomFullName,
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
		RandomParagraph,
		RandomFormat,
		RandomPhone,
		RandomInString,
	}, valueType)
}
