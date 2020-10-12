
### ref: Manage references for large writing projects

By Andrew Stewart ([http://andrewinfosec.com](http://andrewinfosec.com))

`ref` is a CLI tool for managing references in large writing projects. The program uses the filesystem as a database, with one directory per reference number.  

NB: This program assumes macOS.

#### Workflow

Set the location of the database with `$ ref loc`. The location of the database is stored in `~/.ref`. This step does not have to be repeated.

During the writing process, when the time comes to add a new reference, use `$ ref add`. This creates a reference number, copies it to the clipboard, and opens the corresponding directory. You can then copy the files associated with the reference into the directory, and `command+v` the reference number into your manuscript.

Reference numbers increase monotonically, so reference number 1 is followed by reference number 2, and so on. Reference numbers are never reused.

To view the files associated with a reference number us `$ ref <number>`. This will open the `*.html` and `*.pdf` files from the directory associated with the reference number.
