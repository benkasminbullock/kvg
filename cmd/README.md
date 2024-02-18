#FILES IN THIS DIRECTORY

* ''kvg-mode.el'' provides an Emacs editing mode which automatically
renumbers all the XML elements for consistency, and indents the buffer
each time the file is saved (C-x C-s). It requires the user already
has go-mode.el installed. It also uses a hard-coded path for renumber,
so it will require end-user editing to be used correctly.

* ''Makefile'' builds the two Go binaries.

* ''read-write-test.go'' provides a utility which reads and then
writes back out all the files of kvg, and prints a report on which
files differ from the standard formatting.

* ''renumber.go'' provides a utility which reformats and renumbers the
files provided on the command line. This is used by the Emacs editing
mode.
