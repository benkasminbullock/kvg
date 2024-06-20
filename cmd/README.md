# FILES IN THIS DIRECTORY

* __bogusgroup.go__ is a tool to find groups with no paths in them

* __empty-path.go__ finds files where the number of strokes does not
match the number of stroke number labels. It also locates instances
of empty paths with no information. As of 2024-06-20 there are no
instances in the repository.

* __kvg-mode.el__ provides an Emacs editing mode which automatically
renumbers all the XML elements for consistency, and indents the
buffer each time the file is saved (C-x C-s). It requires the user
already has go-mode.el installed. It also uses a hard-coded path for
renumber, so it will require end-user editing to be used correctly.

* __Makefile__ builds the two Go binaries.

* __read-write-test.go__ provides a utility which reads and then
writes back out all the files of kvg, and prints a report on which
files differ from the standard formatting.

* __renumber.go__ provides a utility which reformats and renumbers the
files provided on the command line. This is used by the Emacs editing
mode.

