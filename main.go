// A split diff with large fillers, the hatched blanks across from added
// lines, scrolling by itself for 20 seconds, then saying how many frames it
// drew. Slow frames go to stderr:
//
//	MYGO_FRAME_STATS=1 go run .           # with fillers
//	MYGO_FRAME_STATS=1 go run . -unified  # without, to compare
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// texts returns two texts: the new one adds runs of lines to the old one,
// of different sizes, so that fillers keep coming into view and leaving it.
func texts() (before, after string) {
	var o, n strings.Builder
	for i, run := range []int{30, 90, 15, 60, 120, 25, 75, 40, 100, 20, 50, 80} {
		for j := range 12 {
			line := fmt.Sprintf("unchanged line %d of block %d\n", j, i)
			o.WriteString(line)
			n.WriteString(line)
		}
		for j := range run {
			fmt.Fprintf(&n, "added line %d of block %d, long enough to look like code\n", j, i)
		}
	}
	return o.String(), n.String()
}

func main() {
	unified := flag.Bool("unified", false, "show the diff unified, which has no fillers")
	flag.Parse()

	before, after := texts()
	d := &editor.DiffState{Old: editor.NewBuffer(before), New: editor.NewBuffer(after)}
	if *unified {
		d.Layout = editor.DiffUnified
	}
	// 7.25 DIPs a frame, down and back up, as a trackpad would.
	step, last := 7.25, -1.0
	frames, moved := 0, 0.0
	view := func(c *ui.Context) {
		y := d.ScrollY()
		if y == last {
			step = -step
		}
		frames++
		moved += math.Abs(y - last)
		last = y
		d.SetScrollY(y + step)
		c.Invalidate()
		editor.DiffView(c, d).Fill()
	}

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "Filler scroll", Width: 1600, Height: 1000, Content: ui.View(view)})
		time.AfterFunc(20*time.Second, mygo.App.Quit)
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d frames in 20 seconds, %.0f a second, scrolling %.0f DIPs\n", frames, float64(frames)/20, moved)
}
