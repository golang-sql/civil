// Copyright 2016 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package civil implements types for civil time, a time-zone-independent
// representation of time that follows the rules of the proleptic
// Gregorian calendar with exactly 24-hour days, 60-minute hours, and 60-second
// minutes.
//
// Because they lack location information, these types do not represent unique
// moments or intervals of time. Use time.Time for that purpose.
package civil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate_MarshalJSON(t *testing.T) {

	dLeap := Date{
		Year:  2020,
		Month: 2,
		Day:   29,
	}

	json, err := dLeap.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2020-02-29"`), json)

	dInvalid := Date{
		Year:  -1,
		Month: 2,
		Day:   29,
	}

	json, err = dInvalid.MarshalJSON()
	assert.EqualError(t, err, "Date.MarshalJSON: year '-1' outside of range [0,9999]")
	assert.Nil(t, json)
}

func TestDate_UnmarshalJSON(t *testing.T) {
	jsonLeap := []byte(`"2020-02-29"`)
	dLeap := &Date{}
	err := dLeap.UnmarshalJSON(jsonLeap)
	assert.NoError(t, err)
	assert.Equal(t, Date{Year: 2020, Month: 2, Day: 29}, *dLeap)

	jsonInvalid := []byte(`"2020-13-40"`)
	dInvalid := &Date{}
	err = dInvalid.UnmarshalJSON(jsonInvalid)
	assert.EqualError(t, err, "invalid date: parsing time \"2020-13-40\": month out of range")
}

func TestDate_AddMonths(t *testing.T) {
	dLeap := Date{
		Year:  2020,
		Month: 2,
		Day:   29,
	}

	dLeap = dLeap.AddMonths(12)

	assert.Equal(t, Date{Year: 2021, Month: 3, Day: 1}, dLeap) // no leap day in 2021, so pushes over to 3/1
}

func TestDate_AddYears(t *testing.T) {
	dLeap := Date{
		Year:  2020,
		Month: 2,
		Day:   29,
	}

	dLeap = dLeap.AddYears(1)

	assert.Equal(t, Date{Year: 2021, Month: 3, Day: 1}, dLeap) // no leap day in 2021, so pushes over to 3/1
}

func TestDate_Value(t *testing.T) {
	d := Date{
		Year:  2020,
		Month: 2,
		Day:   29,
	}

	v, err := d.Value()
	assert.NoError(t, err)
	assert.Equal(t, v, "2020-02-29")
}

func TestDate_Scan(t *testing.T) {
	d := &Date{}
	var v interface{}
	v = "2020-02-29"
	d.Scan(v)
	assert.Equal(t, Date{Year: 2020, Month: 2, Day: 29}, *d)
}

func TestTime_MarshalJSON(t *testing.T) {

	time := Time{
		Hour:       3,
		Minute:     42,
		Second:     31,
		Nanosecond: 876,
	}

	json, err := time.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"03:42:31.000000876"`), json)
}

func TestTime_UnmarshalJSON(t *testing.T) {
	jsonGood := []byte(`"03:42:31.000000876"`)
	timeGood := &Time{}
	err := timeGood.UnmarshalJSON(jsonGood)
	assert.NoError(t, err)
	assert.Equal(t, Time{Hour: 3, Minute: 42, Second: 31, Nanosecond: 876}, *timeGood)

	jsonInvalid := []byte(`"-3:42:31.000000876"`)
	timeInvalid := &Time{}
	err = timeInvalid.UnmarshalJSON(jsonInvalid)
	assert.EqualError(t, err, "invalid time: parsing time \"-3:42:31.000000876\" as \"15:04:05.999999999\": cannot parse \"-3:42:31.000000876\" as \"15\"")
}

func TestTime_Value(t *testing.T) {
	time := Time{
		Hour:       3,
		Minute:     42,
		Second:     31,
		Nanosecond: 876,
	}

	v, err := time.Value()
	assert.NoError(t, err)
	assert.Equal(t, "03:42:31.000000876", v)
}

func TestTime_Scan(t *testing.T) {
	time := &Time{}
	var v interface{}
	v = "03:42:31.000000876"
	time.Scan(v)
	assert.Equal(t, *time, Time{Hour: 3, Minute: 42, Second: 31, Nanosecond: 876})
}

func TestDateTime_MarshalJSON(t *testing.T) {

	datetime := DateTime{
		Date: Date{
			Year:  2020,
			Month: 2,
			Day:   29,
		},
		Time: Time{
			Hour:       3,
			Minute:     42,
			Second:     31,
			Nanosecond: 876,
		},
	}

	json, err := datetime.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2020-02-29T03:42:31.000000876"`), json)
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	jsonGood := []byte(`"2020-02-29T03:42:31.000000876"`)
	datetimeGood := &DateTime{}
	err := datetimeGood.UnmarshalJSON(jsonGood)
	assert.NoError(t, err)
	expectedGood := DateTime{
		Date: Date{
			Year:  2020,
			Month: 2,
			Day:   29,
		},
		Time: Time{
			Hour:       3,
			Minute:     42,
			Second:     31,
			Nanosecond: 876,
		},
	}
	assert.Equal(t, expectedGood, *datetimeGood)

	jsonInvalid := []byte(`"0-02-29T03:42:31.000000876"`)
	datetimeInvalid := &DateTime{}
	err = datetimeInvalid.UnmarshalJSON(jsonInvalid)
	assert.EqualError(t, err, "invalid datetime: parsing time \"0-02-29T03:42:31.000000876\" as \"2006-01-02t15:04:05.999999999\": cannot parse \"-29T03:42:31.000000876\" as \"2006\"")
}

func TestDateTime_Value(t *testing.T) {
	datetime := DateTime{
		Date: Date{
			Year:  2020,
			Month: 2,
			Day:   29,
		},
		Time: Time{
			Hour:       3,
			Minute:     42,
			Second:     31,
			Nanosecond: 876,
		},
	}

	v, err := datetime.Value()
	assert.NoError(t, err)
	assert.Equal(t, "2020-02-29T03:42:31.000000876", v)
}

func TestDateTime_Scan(t *testing.T) {
	datetime := &DateTime{}
	var v interface{}
	v = "2020-02-29T03:42:31.000000876"
	datetime.Scan(v)
	expected := DateTime{
		Date: Date{
			Year:  2020,
			Month: 2,
			Day:   29,
		},
		Time: Time{
			Hour:       3,
			Minute:     42,
			Second:     31,
			Nanosecond: 876,
		},
	}
	assert.Equal(t, *datetime, expected)
}
