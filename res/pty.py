#!/usr/bin/python
import pty, sys
pty.spawn(sys.argv[1:])