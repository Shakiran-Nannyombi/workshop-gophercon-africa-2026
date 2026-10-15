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

// flamingoX and flamingoY scatter a flock along the rift lake's shore.
func flamingoX(i int) string { return fmt.Sprint(380 + (i*71)%520) }
func flamingoY(i int) string { return fmt.Sprint(178 + (i*7)%14) }

// stoneTownX and stoneTownY place arched windows on the Stone Town houses.
func stoneTownX(i int) string { return fmt.Sprint([]int{56, 76, 160, 188, 232, 320, 342, 384}[i%8]) }
func stoneTownY(i int) string { return fmt.Sprint([]int{164, 180, 136, 160, 170, 180, 186, 188}[i%8]) }
