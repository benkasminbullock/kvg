This module allows a KanjiVG file to easily be read, parsed, and
written. It also contains functions which enable bulk editing of the
KanjiVG files.

This library automatically renumbers XML ids when writing a file. For
example, if you exchange the order of the text elements in the Go
structure, the values of `id` of the `text` elements in the output
file will be numbered in the order of the Go structure, disregarding
the order they were in originally. If you change the top level name of
the character in the base element, the names of all of the other
elements will be changed. The `id` values of paths and groups are
renumbered in the order they appear in the structure, discarding their
original numbering. In other words, when reordering elements or adding
elements, the library user only needs to think about the tree of
`Child`, `Group`, `Text`, and `Path` elements, and does not need to
give the elements numbers by himself.
