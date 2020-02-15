#!/usr/bin/env python
from __future__ import print_function

import subprocess
import sys


def main():
    try:
        args = ["gofmt", "-l", "-w"]
        args.extend(sys.argv[1:])
        output = subprocess.check_output(args, stderr=subprocess.STDOUT, universal_newlines=True)
        output = output.strip()
        if output:
            print(output)
            return sys.exit(1)
    except subprocess.CalledProcessError as exc:
        print(exc.output.strip())
        return exc.returncode
    except Exception:
        pass
    return sys.exit(0)
