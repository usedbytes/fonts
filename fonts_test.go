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
    "fmt"
    "testing"
)

import _ "image/png"

func TestWrap(t *testing.T) {
    wrapwidth := 101

    font := NewFontFromFile("tiny_font.fnt")
    fmt.Printf("Font name: %s, Number of Characters: %v, Average Width: %v, Height: %v\n", font.name, len(font.letters), font.avg_width, font.height)
    lines := font.WrapText("Today I learned that writing word wrapping for the second time was still around 100,000,000 times harder than I was expecting", wrapwidth)
    for _, l := range lines {
        fmt.Printf("%v\t%v\n", l, font.Width(l))
    }

}
