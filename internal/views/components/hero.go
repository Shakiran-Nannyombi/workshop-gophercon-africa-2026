package components

import "fmt"

// windowX and windowY scatter lit windows across the skyline towers.
func windowX(i int) string { return fmt.Sprint(306 + (i*53)%520) }
func windowY(i int) string { return fmt.Sprint(140 + (i*37)%90) }

// teaRow draws one curved row of tea bushes across the highland ridge.
func teaRow(i int) string {
	y := 160 + i*8
	return fmt.Sprintf("M0 %d C300 %d 700 %d 1200 %d v3 C700 %d 300 %d 0 %d z", y, y-14, y+10, y-6, y+13, y-11, y+3)
}
