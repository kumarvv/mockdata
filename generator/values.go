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

func generateValue(ctx context.Context, column *models.ConfigColumn) (interface{}, error) {
	valueType := column.Type
	gender := 0

	if valueType == valuetypes.SQL {
		// TODO
	} else if valueType == valuetypes.String {
		return column.Value, nil
	} else if valueType == valuetypes.Integer {
		return utils.ToInt(column.Value), nil
	} else if valueType == valuetypes.Float {
		return utils.ToFloat(column.Value), nil
	} else if valueType == valuetypes.Boolean {
		return utils.ToBool(column.Value), nil
	} else if valueType == valuetypes.Date {
		return utils.ToTime(column.Value)
	} else if valueType == valuetypes.DateTime {
		return utils.ToTime(column.Value)
	} else if valueType == valuetypes.UUID {
		return uuid.New().String(), nil
	} else if valueType == valuetypes.RandomString {
		return withCase(ctx, column, withLen(ctx, column, randomdata.SillyName())), nil
	} else if valueType == valuetypes.RandomTitle {
		return withCase(ctx, column, withLen(ctx, column, randomdata.Title(gender))), nil
	} else if valueType == valuetypes.RandomGender {
		return randomdata.SillyName(), nil
	} else if valueType == valuetypes.RandomFirstName {
		return randomdata.FirstName(gender), nil
	} else if valueType == valuetypes.RandomLastName {
		return randomdata.LastName(), nil
	} else if valueType == valuetypes.RandomFullName {
		return randomdata.FullName(gender), nil
	} else if valueType == valuetypes.RandomEmail {
		return randomdata.Email(), nil
	} else if valueType == valuetypes.RandomCurrency {
		return randomdata.Currency(), nil
	} else if valueType == valuetypes.RandomAddress {
		return randomdata.Address(), nil
	} else if valueType == valuetypes.RandomStreet {
		return randomdata.Street(), nil
	} else if valueType == valuetypes.RandomCity {
		return randomdata.City(), nil
	} else if valueType == valuetypes.RandomState {
		return randomdata.State(randomdata.Large), nil
	} else if valueType == valuetypes.RandomState2 {
		return randomdata.State(randomdata.Small), nil
	} else if valueType == valuetypes.RandomCountry {
		return randomdata.Country(randomdata.FullCountry), nil
	} else if valueType == valuetypes.RandomCountry2 {
		return randomdata.Country(randomdata.TwoCharCountry), nil
	} else if valueType == valuetypes.RandomCountry3 {
		return randomdata.Country(randomdata.ThreeCharCountry), nil
	} else if valueType == valuetypes.RandomNumber {
		return randomdata.Number(), nil
	} else if valueType == valuetypes.RandomDecimal {
		return randomdata.Decimal(), nil
	} else if valueType == valuetypes.RandomBoolean {
		return randomdata.Boolean(), nil
	} else if valueType == valuetypes.RandomParagraph {
		return randomdata.Paragraph(), nil
	} else if valueType == valuetypes.RandomFormat {
		return randomdata.StringNumber(3, "-"), nil
	} else if valueType == valuetypes.RandomDate {
		return randomdata.FullDate(), nil
	} else if valueType == valuetypes.RandomDay {
		return randomdata.Day(), nil
	} else if valueType == valuetypes.RandomMonth {
		return randomdata.Month(), nil
	} else if valueType == valuetypes.RandomYear {
		return randomdata.Number(1900, 2999), nil
	} else if valueType == valuetypes.RandomPhone {
		return randomdata.PhoneNumber(), nil
	} else if valueType == valuetypes.RandomIn {
		return randomdata.StringSample("a", "b"), nil
	} else if valueType == valuetypes.RandomFrom {
		// TODO sql
	}

	return nil, nil
}

func withCase(ctx context.Context, column *models.ConfigColumn, value string) string {
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

func withLen(ctx context.Context, column *models.ConfigColumn, value string) string {
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
