### [Files](https://build-your-own.org/database/01_files)

This section basically about how data should be stored in database. It is different from being stored like normal local files.

Key components:

* Atomic (when write data, it must be fully written or nothing, no half-written data)
* Can handle powerloss, unexpected shutdown while writing data
