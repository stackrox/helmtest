# rox-imported directory

The Go packages in subdirectories are stripped-down versions imported from another repository.

They are in the `internal/` directory as the intention is to open-source them as a standalone
library at a later point in time. Once that happens, they can be removed and said library should
be added as a dependency to this repository. In the meantime, however, we do not want anyone to
start pulling them from this repo.
