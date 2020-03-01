# `/pkg`

`/pkg` directory contains library code that's ok to share across services.

It's also a way to group Go code in one place when your root directory contains lots of non-Go components and directories making it easier to run various Go tools.

Some popular go projects use `/pkg` to export library code that is ok to use by external apps but in our case, we use `/pkg` to contains library code that is ok to share across services.
