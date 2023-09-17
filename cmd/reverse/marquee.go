package reverse

import (
	"github.com/mattn/go-runewidth"
	"github.com/vbauerster/mpb/v8/decor"
)

// Marquee returns marquee decorator that will scroll text right to left.
//
// `text` is the scrolling message
//
// `ws` controls the showing window size
//
//	func Marquee(text string, ws uint, wcc ...decor.WC) decor.Decorator {
//		bytes := []byte(text)
//		buf := make([]byte, ws)
//		var count uint
//		f := func(s decor.Statistics) string {
//			start := count % uint(len(bytes))
//			var i uint = 0
//			for ; i < ws && start+i < uint(len(bytes)); i++ {
//				buf[i] = bytes[start+i]
//			}
//			for ; i < ws; i++ {
//				buf[i] = ' '
//			}
//			count++
//			return string(buf)
//		}
//		return decor.Any(f, wcc...)
//	}
func Marquee(t string, ws int, divider string, wcc ...decor.WC) decor.Decorator {

	var count int
	var f = func(s decor.Statistics) string {
		length := runewidth.StringWidth(t + divider)
		if runewidth.StringWidth(t) <= ws {
			return runewidth.FillRight(t, ws)
		}
		text := t + divider
		var msg string

		msg = TruncateLeft(Truncate(text, count+ws), count)
		if ws+count > length {
			// start         end       ws
			// text[count:ws+count]
			// msg = TruncateLeft(text, length-count)
			// start         end       ws
			// text[0:ws-count-length]
			// msg += fmt.Sprintf("%d", count+ws-length)
			msg += Truncate(text, count+ws-length)
		}
		count++
		if count >= length {
			count = 0
		}
		return runewidth.FillRight(msg, ws)

	}
	return decor.Any(f, wcc...)
}

func Truncate(s string, size int) string {
	return runewidth.Truncate(s, size, "")
}
func TruncateLeft(s string, size int) string {
	return runewidth.TruncateLeft(s, size, "")
}
