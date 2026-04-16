package generator

import (
	"context"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"kumarvv.com/mockdata/constants/valuetypes"
	"kumarvv.com/mockdata/models"
	"kumarvv.com/mockdata/utils"
)

type valueGen struct {
	column *models.ConfigColumn
	value  interface{}
	err    error
}

func getValue(column *models.ConfigColumn, value interface{}) *valueGen {
	return &valueGen{
		column: column,
		value:  value,
		err:    nil,
	}
}
func (v *valueGen) Value() (interface{}, error) {
	if valuetypes.IsString(v.column.Type) {
		valueStr := utils.ToString(v.value)
		valueStr = withLen(v.column, valueStr)
		valueStr = withCase(v.column, valueStr)
		return valueStr, v.err
	} else {
		return v.value, v.err
	}
}

func generateValue(ctx context.Context, column *models.ConfigColumn, gender, ix int) (interface{}, error) {
	valueType := column.Type
	value := column.Value
	var err error
	if valueType == valuetypes.SQL {
		// TODO
	} else if valueType == valuetypes.String {
		value, err = getValue(column, utils.ToString(column.Value)).Value()
	} else if valueType == valuetypes.Integer {
		value, err = getValue(column, utils.ToInt(column.Value)).Value()
	} else if valueType == valuetypes.Float {
		value, err = getValue(column, utils.ToFloat(column.Value)).Value()
	} else if valueType == valuetypes.Boolean {
		value, err = getValue(column, utils.ToBool(column.Value)).Value()
	} else if valueType == valuetypes.Date {
		if column.Format != nil {
			value, err = utils.ToTimeFormat(column.Value, *column.Format)
		} else {
			value, err = utils.ToTime(column.Value)
		}
	} else if valueType == valuetypes.DateTime {
		if column.Format != nil {
			value, err = utils.ToTimeFormat(column.Value, *column.Format)
		} else {
			value, err = utils.ToTime(column.Value)
		}
	} else if valueType == valuetypes.Serial {
		minValue := 1
		if column.Min != nil {
			minValue = *column.Min
		}
		value = minValue + ix
	} else if valueType == valuetypes.UUID {
		value = uuid.New().String()
	} else if valueType == valuetypes.RandomString {
		value, err = getValue(column, randomdata.SillyName()).Value()
	} else if valueType == valuetypes.RandomTitle {
		value, err = getValue(column, randomdata.Title(gender)).Value()
	} else if valueType == valuetypes.RandomGender {
		valueStr := "male"
		if gender == randomdata.Female {
			valueStr = "female"
		}
		value, err = getValue(column, valueStr).Value()
	} else if valueType == valuetypes.RandomFirstName {
		value, err = getValue(column, randomdata.FirstName(gender)).Value()
	} else if valueType == valuetypes.RandomLastName {
		value, err = getValue(column, randomdata.LastName()).Value()
	} else if valueType == valuetypes.RandomFullName {
		value, err = getValue(column, randomdata.FullName(gender)).Value()
	} else if valueType == valuetypes.RandomEmail {
		value, err = getValue(column, randomdata.Email()).Value()
	} else if valueType == valuetypes.RandomCurrency {
		value, err = getValue(column, randomdata.Currency()).Value()
	} else if valueType == valuetypes.RandomAddress {
		value, err = getValue(column, randomdata.Address()).Value()
	} else if valueType == valuetypes.RandomStreet {
		value, err = getValue(column, randomdata.Street()).Value()
	} else if valueType == valuetypes.RandomCity {
		value, err = getValue(column, randomdata.City()).Value()
	} else if valueType == valuetypes.RandomState {
		value, err = getValue(column, randomdata.State(randomdata.Large)).Value()
	} else if valueType == valuetypes.RandomState2 {
		value, err = getValue(column, randomdata.State(randomdata.Small)).Value()
	} else if valueType == valuetypes.RandomCountry {
		value, err = getValue(column, randomdata.Country(randomdata.FullCountry)).Value()
	} else if valueType == valuetypes.RandomCountry2 {
		value, err = getValue(column, randomdata.Country(randomdata.TwoCharCountry)).Value()
	} else if valueType == valuetypes.RandomCountry3 {
		value, err = getValue(column, randomdata.Country(randomdata.ThreeCharCountry)).Value()
	} else if valueType == valuetypes.RandomNumber {
		if column.Min != nil && column.Max != nil {
			value, err = getValue(column, randomdata.Number(*column.Min, *column.Max)).Value()
		} else if column.Min != nil {
			value, err = getValue(column, randomdata.Number(*column.Min)).Value()
		} else if column.Max != nil {
			value, err = getValue(column, randomdata.Number(0, *column.Max)).Value()
		} else {
			value, err = getValue(column, randomdata.Number()).Value()
		}
	} else if valueType == valuetypes.RandomDecimal {
		if column.Min != nil && column.Max != nil {
			value, err = getValue(column, randomdata.Decimal(*column.Min, *column.Max)).Value()
		} else if column.Min != nil {
			value, err = getValue(column, randomdata.Decimal(*column.Min)).Value()
		} else if column.Max != nil {
			value, err = getValue(column, randomdata.Decimal(0, *column.Max)).Value()
		} else {
			value, err = getValue(column, randomdata.Decimal()).Value()
		}
	} else if valueType == valuetypes.RandomBoolean {
		value, err = getValue(column, randomdata.Boolean()).Value()
	} else if valueType == valuetypes.RandomParagraph {
		value, err = getValue(column, randomdata.Paragraph()).Value()
	} else if valueType == valuetypes.RandomFormat {
		value, err = getValue(column, randomdata.Boolean()).Value()
		if column.NumPairs != nil && column.Separator != nil {
			return randomdata.StringNumber(*column.NumPairs, *column.Separator), nil
		} else {
			randomdata.StringNumber(1, "")
		}
	} else if valueType == valuetypes.RandomDate {
		dtStr := randomdata.FullDate()
		dt, err := utils.ToTimeFormat(dtStr, randomdata.DateOutputLayout)
		if err == nil {
			if column.Format != nil {
				value = dt.Format(*column.Format)
			} else {
				value = dt.Format(utils.DateFormatYMD)
			}
		}
	} else if valueType == valuetypes.RandomDay {
		value = randomdata.Day()
	} else if valueType == valuetypes.RandomMonth {
		value = randomdata.Month()
	} else if valueType == valuetypes.RandomYear {
		value = randomdata.Number(1900, 2999)
	} else if valueType == valuetypes.RandomPhone {
		value = randomdata.PhoneNumber()
	} else if valueType == valuetypes.RandomInString {
		valueStr := utils.ToString(column.Value)
		tokens := strings.Split(valueStr, ",")
		value = utils.RandomOneOf(tokens...)
	} else if valueType == valuetypes.RandomInInteger {
		valueStr := utils.ToString(column.Value)
		tokens := utils.SplitToInt(valueStr, ",")
		value = utils.RandomOneOf(tokens...)
	} else if valueType == valuetypes.RandomInFloat {
		valueStr := utils.ToString(column.Value)
		tokens := utils.SplitToFloat(valueStr, ",")
		value = utils.RandomOneOf(tokens...)
	} else if valueType == valuetypes.RandomFromSQL {
		// TODO sql
	}

	return value, err
}

func withCase(column *models.ConfigColumn, value string) string {
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

func withLen(column *models.ConfigColumn, value string) string {
	if column.Max == nil && column.Min == nil {
		return value
	}

	lValue := value
	if column.Min != nil && len(value) < *column.Min {
		for len(lValue) < *column.Min {
			lValue = value + randomdata.SillyName()
		}
	}
	if len(lValue) > *column.Max {
		return lValue[:*column.Max]
	}

	return lValue
}
