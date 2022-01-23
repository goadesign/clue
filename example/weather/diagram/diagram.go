package model

import . "goa.design/model/dsl"

var _ = Design("Weather System Architecture", "The Weather example system architecture", func() {
	var WeatherGov = SoftwareSystem("weather.gov", "Provides weather forecasts for US locations.", func() {
		External()
		URL("https://weather.go")
		Tag("External")
	})

	var IPAPI = SoftwareSystem("ip-api.com", "Provides IP geolocation information.", func() {
		External()
		URL("https://ip-api.com")
		Tag("External")
	})

	var _ = Person("User", "A client of the weather system.", func() {
		Uses("Weather Software System/Front Service", "Makes requests to", "HTTP", Synchronous)
	})

	var System = SoftwareSystem("Weather Software System", "Provides IP based weather forecasts.", func() {
		var Location = Container("Location Service", "Leverages ip-api.com to locate IP addresses.", "Go and Goa", func() {
			Uses(IPAPI, "Makes requests to", "HTTP", Synchronous)
			Tag("Service")
		})

		var Forecast = Container("Forecast Service", "Leverages weather.gov to retrieve weather forecasts for US based locations.", "Go and Goa", func() {
			Uses(WeatherGov, "Makes requests to", "HTTP", Synchronous)
			Tag("Service")
		})

		Container("Front Service", "Retrieves weather forecasts for a given IP.", "Go and Goa", func() {
			Tag("Service")
			Uses(Location, "Retrieves IP location from", "gRPC", Synchronous)
			Uses(Forecast, "Retrieves weather forecasts from", "gRPC", Asynchronous)
		})
	})

	Views(func() {
		ContainerView(System, "Weather System Services", "Weather software system architecture diagram", func() {
			AddAll()
			AutoLayout(RankLeftRight)
		})
		Styles(func() {
			ElementStyle("Person", func() {
				Background("#e6e6ea")
				Stroke("#f75c03")
				Shape(ShapePerson)
			})
			ElementStyle("Container", func() {
				Background("#e6e6ea")
				Stroke("#2ab7ca")
			})
			ElementStyle("External", func() {
				Background("#eae6e6")
				Stroke("#cab72a")
			})
			ElementStyle("Software System", func() {
				Shape(ShapeRoundedBox)
				Background("#e6e6ea")
				Stroke("#f75c03")
			})
		})
	})
})
