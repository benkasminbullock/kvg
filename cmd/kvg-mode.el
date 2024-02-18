;;; kvg-mode.el --- Perl code editing commands   -*- lexical-binding:t -*-

(require 'go-mode)

(defcustom renumber-command "/home/ben/projects/kvgpub/cmd/renumber"
  "The renumber command"
  :type 'string
  :group 'kvg)
(setq rng-nxml-auto-validate-flag nil)
(defvar kvg-mode-font-lock-keywords
  '(("\\<kvg\\>" . font-lock-keyword-face)))
(define-derived-mode kvg-mode nxml-mode "Kvg"
  "A major mode to edit GNU ld script files."
  (font-lock-add-keywords nil kvg-mode-font-lock-keywords)
  (setq tab-width 4
	indent-tabs-mode t
	nxml-child-indent 4
	nxml-attribute-indent 4)
  (add-hook 'before-save-hook 'renumber-before-save))

(defun renumber-before-save ()
  "Run renumber on the file before writing it"
  (interactive)
  (when (eq major-mode 'kvg-mode) (renumber)))

(defun renumber ()
  "Renumber the groups and reindent using the tool"
  (interactive)
  (let ((tmpfile (make-temp-file "renumber" nil ".svg"))
	(patchbuf (get-buffer-create "*Renumber patch*"))
	(errbuf (get-buffer-create "*Renumber Errors*"))
        (coding-system-for-read 'utf-8)
        (coding-system-for-write 'utf-8))
    (save-restriction
      (widen)
      (if errbuf
	  (with-current-buffer errbuf
	    (setq buffer-read-only nil)
	    (erase-buffer)))
      (with-current-buffer patchbuf
	    (erase-buffer))
      (write-region nil nil tmpfile)
      (if (zerop (call-process renumber-command nil errbuf nil tmpfile))
	  (progn
	    (if (zerop (call-process-region (point-min) (point-max) "diff" nil patchbuf nil "-n" "-" tmpfile))
		(message "Buffer format already OK!")
	      (go--apply-rcs-patch patchbuf)
	      (message "Applied renumber"))
	    (if errbuf (gofmt--kill-error-buffer errbuf)))
	(message "Could not apply renumber")
	(if errbuf (gofmt--process-errors (buffer-file-name) tmpfile errbuf)))
      (kill-buffer patchbuf)
      (delete-file tmpfile))))

(provide 'kvg-mode)
