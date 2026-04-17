package generator

import (
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"kumarvv.com/mockdata/constants/functiontypes"
	"kumarvv.com/mockdata/models"
	"kumarvv.com/mockdata/utils"
)

type valueGen struct {
	column *models.Column
	value  interface{}
	err    error
}

func getValue(column *models.Column, value interface{}) *valueGen {
	return &valueGen{
		column: column,
		value:  value,
		err:    nil,
	}
}
func (v *valueGen) Value() (interface{}, error) {
	if functiontypes.IsString(v.column.FnName) {
		valueStr := utils.ToString(v.value)
		valueStr = withLen(v.column, valueStr)
		valueStr = withCase(v.column, valueStr)
		return valueStr, v.err
	} else {
		return v.value, v.err
	}
}

func generateValue(table *models.ConfigTable, column *models.Column, gender, ix int) (interface{}, error) {
	valueType := column.FnName
	value := column.Value
	var err error
	if valueType == functiontypes.SQL {
		// TODO
	} else if valueType == functiontypes.String {
		value, err = getValue(column, utils.ToString(column.Value)).Value()
	} else if valueType == functiontypes.Integer {
		value, err = getValue(column, utils.ToInt64(column.Value)).Value()
	} else if valueType == functiontypes.Float {
		value, err = getValue(column, utils.ToFloat(column.Value)).Value()
	} else if valueType == functiontypes.Boolean {
		value, err = getValue(column, utils.ToBool(column.Value)).Value()
	} else if valueType == functiontypes.Date {
		if column.Format != nil {
			value, err = utils.ToTimeFormat(column.Value, *column.Format)
		} else {
			value, err = utils.ToTime(column.Value)
		}
	} else if valueType == functiontypes.DateTime {
		if column.Format != nil {
			value, err = utils.ToTimeFormat(column.Value, *column.Format)
		} else {
			value, err = utils.ToTime(column.Value)
		}
	} else if valueType == functiontypes.Serial {
		minValue := 1
		if column.Min != nil {
			minValue = *column.Min
		}
		value = minValue + ix
	} else if valueType == functiontypes.UUID {
		value = uuid.New().String()
	} else if valueType == functiontypes.RandomString {
		value, err = getValue(column, randomdata.SillyName()).Value()
	} else if valueType == functiontypes.RandomTitle {
		value, err = getValue(column, randomdata.Title(gender)).Value()
	} else if valueType == functiontypes.RandomGender {
		valueStr := "male"
		if gender == randomdata.Female {
			valueStr = "female"
		}
		value, err = getValue(column, valueStr).Value()
	} else if valueType == functiontypes.RandomFirstName {
		value, err = getValue(column, randomdata.FirstName(gender)).Value()
	} else if valueType == functiontypes.RandomLastName {
		value, err = getValue(column, randomdata.LastName()).Value()
	} else if valueType == functiontypes.RandomFullName {
		value, err = getValue(column, randomdata.FullName(gender)).Value()
	} else if valueType == functiontypes.RandomEmail {
		value, err = getValue(column, randomdata.Email()).Value()
	} else if valueType == functiontypes.RandomCurrency {
		value, err = getValue(column, randomdata.Currency()).Value()
	} else if valueType == functiontypes.RandomAddress {
		value, err = getValue(column, randomdata.Address()).Value()
	} else if valueType == functiontypes.RandomStreet {
		value, err = getValue(column, randomdata.Street()).Value()
	} else if valueType == functiontypes.RandomCity {
		value, err = getValue(column, randomdata.City()).Value()
	} else if valueType == functiontypes.RandomState {
		value, err = getValue(column, randomdata.State(randomdata.Large)).Value()
	} else if valueType == functiontypes.RandomState2 {
		value, err = getValue(column, randomdata.State(randomdata.Small)).Value()
	} else if valueType == functiontypes.RandomCountry {
		value, err = getValue(column, randomdata.Country(randomdata.FullCountry)).Value()
	} else if valueType == functiontypes.RandomCountry2 {
		value, err = getValue(column, randomdata.Country(randomdata.TwoCharCountry)).Value()
	} else if valueType == functiontypes.RandomCountry3 {
		value, err = getValue(column, randomdata.Country(randomdata.ThreeCharCountry)).Value()
	} else if valueType == functiontypes.RandomNumber {
		if column.Min != nil && column.Max != nil {
			value, err = getValue(column, randomdata.Number(*column.Min, *column.Max)).Value()
		} else if column.Min != nil {
			value, err = getValue(column, randomdata.Number(*column.Min)).Value()
		} else if column.Max != nil {
			value, err = getValue(column, randomdata.Number(0, *column.Max)).Value()
		} else {
			value, err = getValue(column, randomdata.Number()).Value()
		}
	} else if valueType == functiontypes.RandomDecimal {
		if column.Min != nil && column.Max != nil {
			value, err = getValue(column, randomdata.Decimal(*column.Min, *column.Max)).Value()
		} else if column.Min != nil {
			value, err = getValue(column, randomdata.Decimal(*column.Min)).Value()
		} else if column.Max != nil {
			value, err = getValue(column, randomdata.Decimal(0, *column.Max)).Value()
		} else {
			value, err = getValue(column, randomdata.Decimal()).Value()
		}
	} else if valueType == functiontypes.RandomBoolean {
		value, err = getValue(column, randomdata.Boolean()).Value()
	} else if valueType == functiontypes.RandomParagraph {
		value, err = getValue(column, randomdata.Paragraph()).Value()
	} else if valueType == functiontypes.RandomFormat {
		value, err = getValue(column, randomdata.Boolean()).Value()
		if column.NumPairs != nil && column.Separator != nil {
			return randomdata.StringNumber(*column.NumPairs, *column.Separator), nil
		} else {
			randomdata.StringNumber(1, "")
		}
	} else if valueType == functiontypes.RandomDate {
		dtStr := randomdata.FullDate()
		dt, err := utils.ToTimeFormat(dtStr, randomdata.DateOutputLayout)
		if err == nil {
			if column.Format != nil {
				value = dt.Format(*column.Format)
			} else {
				value = dt.Format(utils.DateFormatYMD)
			}
		}
	} else if valueType == functiontypes.RandomDay {
		value = randomdata.Day()
	} else if valueType == functiontypes.RandomMonth {
		value = randomdata.Month()
	} else if valueType == functiontypes.RandomYear {
		value = randomdata.Number(1900, 2999)
	} else if valueType == functiontypes.RandomPhone {
		value = randomdata.PhoneNumber()
	} else if valueType == functiontypes.RandomInString {
		valueStr := utils.ToString(column.Value)
		tokens := strings.Split(valueStr, ",")
		value = utils.RandomOneOf(tokens...)
	} else if valueType == functiontypes.RandomInInteger {
		valueStr := utils.ToString(column.Value)
		tokens := utils.SplitToInt(valueStr, ",")
		value = utils.RandomOneOf(tokens...)
	} else if valueType == functiontypes.RandomInFloat {
		valueStr := utils.ToString(column.Value)
		tokens := utils.SplitToFloat(valueStr, ",")
		value = utils.RandomOneOf(tokens...)
	} else if valueType == functiontypes.RandomRange {
		if column.Min != nil && column.Max != nil {
			value = randomdata.Number(*column.Min, *column.Max)
		} else if column.Min != nil {
			value = randomdata.Number(*column.Min, *column.Min+table.RowCount)
		} else if column.Max != nil {
			value = randomdata.Number(1, *column.Max)
		} else {
			value = randomdata.Number()
		}
	} else if valueType == functiontypes.RandomFromSQL {
		// TODO sql
	}

	return value, err
}

func withCase(column *models.Column, value string) string {
	if column.Case == nil {
		return value
	}
	if *column.Case == "lower" {
		return strings.ToLower(value)
	} else if *column.Case == "upper" {
		return strings.ToUpper(value)
	} else {
		return value
	}
}

func withLen(column *models.Column, value string) string {
	if column.Max == nil && column.Min == nil {
		return value
	}

	lValue := value

	if column.Len != nil {
		for len(lValue) < *column.Len {
			lValue = value + randomdata.SillyName()
		}
		return lValue[:*column.Len]
	}

	if column.Min != nil && len(value) < *column.Min {
		for len(lValue) < *column.Min {
			lValue = value + randomdata.SillyName()
		}
	}
	if column.Max != nil && len(lValue) > *column.Max {
		return lValue[:*column.Max]
	}

	return lValue
}
