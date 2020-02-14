#!/usr/bin/env python
from __future__ import print_function

import subprocess
import sys


def main():
    try:
        args = ['git', 'status', '-u', '--porcelain', '--no-column']
        output = subprocess.check_output(args, stderr=subprocess.STDOUT, universal_newlines=True)
        items = [str(x) for x in output.strip().splitlines()]
        untracked = []
        for item in items:
            if item.startswith('??'):
                untracked.append(item)
        if untracked:
            print('Found untracked files:')
            print('\n'.join(untracked))
            return sys.exit(1)
    except Exception:
        pass
    return sys.exit(0)
