# mygo-filler-scroll

A minimal app showing that a split `editor.DiffView` with large fillers, the hatched blanks across from added lines, drops frames as it scrolls.

It scrolls a generated diff by itself for 20 seconds, down and back up, then prints how many frames it drew and quits. Slow frames go to stderr.

```sh
MYGO_FRAME_STATS=1 go run .           # split, with fillers
MYGO_FRAME_STATS=1 go run . -unified  # unified, without fillers, to compare
```

## Results

On an M3 Mac with a 5K display at 60 Hz, drawn at 4096 × 2304 (looks like 2048 × 1152), 20 seconds each:

| | Frames a second | Frames logged as slow | Over 33 ms | Longest |
|---|---|---|---|---|
| Split, before the fix (MyGo at `7d8b4144`) | 50 | 757 of 1003 | 45 | 48.7 ms |
| Split, after the fix | 60 | 1 of 1195 | 0 | 23.5 ms |
| Unified, without fillers, before the fix | 60 | 1 of 1197 | 0 | 12.8 ms |

It hasn't been tried on other machines, where the difference may be smaller.

To try the fix, from the [`polyline-one-segment`](https://github.com/overcat/mygo/tree/polyline-one-segment) branch:

```sh
git clone -b polyline-one-segment https://github.com/overcat/mygo ../mygo
go mod edit -replace github.com/egoist/mygo=../mygo
```

It needs MyGo after v0.7.4, for `DiffState.SetScrollY`.
