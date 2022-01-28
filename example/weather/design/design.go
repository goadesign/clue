/*
Package design contains shared design resources.
*/
package design

import . "goa.design/goa/v3/dsl"

var Forecast = Type("Forecast", func() {
	Description("Weather forecast")
	Field(1, "location", Location, "Forecast location")
	Field(2, "periods", ArrayOf(Period), "Weather forecast periods")
	Required("location", "periods")
})

var Location = Type("Location", func() {
	Description("Geographical location")
	Field(1, "lat", Float64, "Latitude", func() {
		Example(37.8267)
	})
	Field(2, "long", Float64, "Longitude", func() {
		Example(-122.4233)
	})
	Field(3, "city", String, "City", func() {
		Example("San Francisco")
	})
	Field(4, "state", String, "State", func() {
		Example("CA")
	})
	Required("lat", "long", "city", "state")
})

var Period = Type("Period", func() {
	Description("Weather forecast period")
	Field(1, "name", String, "Period name", func() {
		Example("Morning")
	})
	Field(2, "startTime", String, "Start time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(3, "endTime", String, "End time", func() {
		Format(FormatDateTime)
		Example("2020-01-01T00:00:00Z")
	})
	Field(4, "temperature", Int, "Temperature", func() {
		Example(70)
	})
	Field(5, "temperatureUnit", String, "Temperature unit", func() {
		Example("F")
	})
	Field(6, "summary", String, "Summary", func() {
		Example("Clear")
	})
	Required("name", "startTime", "endTime", "temperature", "temperatureUnit", "summary")
})
