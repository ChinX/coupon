package module

import (
	"math"
)

const earthR float64 = 6378.137

func ToRad(d float64) float64 {
	return d * math.Pi / 180
}

func GetDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radLat1 := ToRad(lat1)
	radLat2 := ToRad(lat2)
	deltaLat := radLat1 - radLat2
	deltaLng := ToRad(lng1) - ToRad(lng2)
	dis := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(deltaLat/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(deltaLng/2), 2)))
	return math.Round(dis*earthR*10000) / 10000
}
