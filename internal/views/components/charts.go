package components

import (
	"fmt"
	"strconv"

	"workshop/internal/data"
)

// Chart geometry shared by the small weather charts (viewBox 360×200).
const (
	chartL, chartR, chartT, chartB = 34.0, 352.0, 14.0, 172.0
)

func colX(i, n int) float64 {
	step := (chartR - chartL) / float64(n)
	return chartL + step*float64(i) + step/2
}

func scaleY(v, lo, hi float64) float64 {
	return chartB - (v-lo)/(hi-lo)*(chartB-chartT)
}

func fx(v float64) string { return strconv.FormatFloat(v, 'f', 1, 64) }

func tempBar(m data.Month, i, n int) (x, y, h string) {
	top, bot := scaleY(m.HighC, 0, 30), scaleY(m.LowC, 0, 30)
	return fx(colX(i, n) - 5), fx(top), fx(bot - top)
}

func rainBar(m data.Month, i, n int) (x, y, h string) {
	top := scaleY(m.RainMM, 0, 200)
	return fx(colX(i, n) - 7), fx(top), fx(chartB - top)
}

func monthX(i, n int) string { return fx(colX(i, n)) }

func gridY(v, lo, hi float64) string { return fx(scaleY(v, lo, hi)) }

func tempTip(m data.Month) string {
	return fmt.Sprintf("%s: ~%v°C high, ~%v°C low", m.Name, m.HighC, m.LowC)
}

func rainTip(m data.Month) string { return fmt.Sprintf("%s: ~%v mm rain", m.Name, m.RainMM) }

func hitX(i, n int) string { return fx(colX(i, n) - (chartR-chartL)/float64(n)/2) }

func hitW(n int) string { return fx((chartR - chartL) / float64(n)) }

func fmtInt(v float64) string { return strconv.Itoa(int(v + 0.5)) }
