/*
 *
 * fonts.go - bitmap font library using 'image'
 *   Copyright Brian Starkey 2013-2014 <stark3y@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of version 2 of the GNU General Public License as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package fonts

import (
    "image"
    "image/draw"
    "image/color"
    "fmt"
    "os"
    "encoding/csv"
    "strconv"
    "strings"
    "unicode"
    "path/filepath"
)

import _ "image/png"

type Font struct {
    name string
    glyphs *image.Alpha
    letters map[byte]image.Rectangle
    avg_width float32
    height int
}

const WRAP_THRESH = 10

func NewFontFromFile(filename string) *Font {

    fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
    defer fp.Close()
	csvr := csv.NewReader(fp)

    header, err := csvr.Read()
    if (err != nil) {
        panic(err)
    } else if (len(header) != 3) {
        panic(fmt.Sprintf("Invalid number of fields in header. " +
            "Expected 3, got %i\n", len(header)))
    }

    name := strings.Trim(header[0], " \t")
    numchars, err := strconv.Atoi(strings.Trim(header[2], " \t"))
    if err != nil {
        panic(err)
    }
    letters := make(map[byte]image.Rectangle, numchars)

    image_file := strings.Trim(header[1], " \t")
    if !filepath.IsAbs(image_file) {
        image_file = filepath.Join(filepath.Dir(filename), image_file)
    }
    fp, err = os.Open(image_file)
    if err != nil {
		panic(err)
	}
    defer fp.Close()

    from_file, _, err := image.Decode(fp)
    if err != nil {
		panic(err)
	}

    records, err := csvr.ReadAll()
    if (err != nil) {
        panic(err)
    }

    for _, l := range records {
        char, _ := strconv.Atoi(strings.Trim(l[0], " \t"))
        left, _ := strconv.Atoi(strings.Trim(l[1], " \t"))
        width, _ := strconv.Atoi(strings.Trim(l[2], " \t"))
        letters[byte(char)] = image.Rect(left, 0, left + width,
            from_file.Bounds().Max.Y)
    }

    return NewFontFromImage(name, &from_file, letters)
}

func NewFontFromImage(name string, glyphs *image.Image,
        letters map[byte]image.Rectangle) *Font {

    font := new(Font)
    font.name = name
    font.letters = letters

    font.avg_width = 0
    for _, v := range letters {
        font.avg_width += float32(v.Dx())
    }
    font.avg_width /= float32(len(letters))

    font.height = (*glyphs).Bounds().Max.Y

    font.glyphs = image.NewAlpha((*glyphs).Bounds())
    draw.Draw(font.glyphs, font.glyphs.Bounds(), *glyphs,
        image.ZP, draw.Src)

    return font
}

func (f *Font) Name() string {
    return f.name
}

func (f *Font) Letters() map[byte]image.Rectangle {
    return f.letters
}

func (f *Font) Width(word string) int {
    length := 0
    for _, l := range word {
        r, ok := f.letters[byte(l)]
        if (!ok) {
            r = f.letters[255]
        }
        length += r.Dx()
    }

    return int(length)
}

func (f *Font) Height() int {
    return f.height
}

// MakeWord makes an image from a string, black on white
func (f *Font) MakeWord(word string) *image.Paletted {

    return f.MakeWordColor(word, color.White, color.Black)
}

// MakeWordColor makes an image from a string (with colors)
func (f *Font) MakeWordColor(word string, bg, fg color.Color) *image.Paletted {
    p := color.Palette{bg, fg}
    r := image.Rect(0, 0, f.Width(word), f.glyphs.Bounds().Dy())

    src := &image.Uniform{fg}

    pal := image.NewPaletted(r, p)

    left := 0
    for _, l := range word {
        r, ok := f.letters[byte(l)]
        if (!ok) {
            r = f.letters[255]
        }

        mask := f.glyphs.SubImage(r)
        draw.DrawMask(pal, image.Rect(left, 0, left + r.Dx(), r.Dy()), 
            src, image.ZP, mask, r.Min, draw.Over)

        left += r.Dx()
    }

    return pal
}

// A slice of s up to the returned value will fit in width
func (f *Font) findSplit(s string, width int) (int) {
    if (f.Width(s) <= width) {
        return len(s)
    }

    n := int(float32(width) / f.avg_width) // approx number of characters

    w := f.Width(s[0:n])


    for ; (w < width) && (n < (len(s) - 1)); w = f.Width(s[0:n]) {
        n++
    }

    for ; (w > width) && (n > 0); w = f.Width(s[0:n]) {
        n--
    }

    return n
}


// WrapText splits a string into chunks which will fit in a given width.
// It will hyphenate in some cases, but numbers will never be split.
func (f *Font) WrapText(sentence string, width int) ([]string) {
    if ( f.Width(sentence) < width) {
        return []string{sentence}
    }

    tokens := strings.Fields(sentence)

    pos := 0
    lines := make([]string, 0, len(tokens))
    thisline := make([]string, 0, len(tokens))

    s_width := f.Width(" ")
    d_width := f.Width("-")

    for i := 0; i < len(tokens); {
        t := tokens[i]
        t_width := f.Width(t)

        end := pos + t_width;
        if (end < width) {
            pos = end + s_width
            thisline = append(thisline, t)
            i++
        } else {
            remaining := width - pos

            if (remaining >= WRAP_THRESH) {
                short_width := remaining - d_width
                split_pos := f.findSplit(t, short_width)
                // Split and hyphenate
                var beginning string
                if (unicode.IsLetter(rune(t[split_pos]))) {
                    beginning = strings.Join([]string{t[0:split_pos], "-"}, "")
                } else {
                    if ( f.Width(string(t[0:split_pos + 1])) <= remaining ) {
                       split_pos++
                    }
                    beginning = t[0:split_pos]
                }

                thisline = append(thisline, beginning)
                tokens[i] = t[split_pos:] // re-evaluate the rest of the word
            }
            pos = 0
        }

        // If the line is complete, or we have reached the end of the sentence
        // put the line together
        if ((pos == 0) || (i == len(tokens))) {
            lines = append(lines[:], strings.Join(thisline, " "))
            thisline = make([]string, 0, len(tokens))
        }
    }

    return lines
}
