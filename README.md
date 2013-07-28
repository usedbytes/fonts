fonts
=====

This is a bitmap font library for go. It can load a font definition and glyphs
from a file, and then uses the image/draw functions to paint them onto images.
It also implements a WrapText function which can be used to split a string
into chunks which will fit into a given width for that font.

Font file format
----------------
A font consists of 2 files - a textual description and an image with the glyphs.
The image should have all the of the glyphs in a single row, appropriately
aligned to put the baseline in the desred position.
The text file consists of a single line header followed by the glyph
definitions.

The header looks like this:

    [font name], [glyph image file], [number of glyphs]

Followed by [number of glyph] lines describing the glyphs:

    [ascii code], [left-most pixel in image], [width of glyph in pixels]

