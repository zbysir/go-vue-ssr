package ssrtool

import (
	"bytes"
)

// 最简单高效的格式化html
func FormatHtml(src string, indent int) (out string) {
	var o bytes.Buffer
	ind := 0
	i := 0
	save := 0
	// <div></div>
	for i < len(src) {
		if src[i] == '<' {
			if src[i+1] == '/' {
				i += 2
				i = skipTo(i, src, '>')
				if i == -1 {
					break
				}
				// 在标签闭合之前添加换行符和空格
				i++
				o.WriteByte('\n')
				o.Write(repeat(' ', ind-2))
				o.WriteString(src[save:i])
				save = i
				ind -= 2
			} else {
				i++
				i = skipTo(i, src, '>')
				if i == -1 {
					break
				}

				i++
				o.WriteByte('\n')
				o.Write(repeat(' ', ind))
				o.WriteString(src[save:i])
				save = i
				ind += 2
			}
		} else {
			i = skipTo(i, src, '<')
			if i == -1 {
				break
			}
			o.WriteByte('\n')
			o.Write(repeat(' ', ind))
			o.WriteString(src[save:i])
			save = i
		}
	}
	o.WriteString(src[save:])

	return o.String()
}

func repeat(x byte, i int) (bs []byte) {
	bs = make([]byte, i)
	for i := range bs {
		bs[i] = x
	}
	return
}

func skipTo(start int, src string, toChar byte) int {
	for start < len(src) {
		if src[start] == toChar {
			return start
		}
		start++
	}
	return -1
}
