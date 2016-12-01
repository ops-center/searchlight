#!/usr/bin/env python

from hyperalert import *
import os
import sys


def call(cmd, stdin=None, cwd=BACKEND_ROOT):
    print(cmd)
    return subprocess.call([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def die(status):
    if status:
        sys.exit(status)


def _go_build(name, goos, goarch):
    print 'Building {name}-{goos}-{goarch}'.format(name=name, goos=goos, goarch=goarch)

    if not os.path.isdir('dist'):
        os.makedirs('dist')

    cmd = "GOOS={goos} GOARCH={goarch} {goc} build -o dist/{name}-{goos}-{goarch} cmd/{name}/*.go".format(
        name=name,
        goc=GOC,
        goos=goos,
        goarch=goarch,
    )
    die(call(cmd))
    print '\n'


def build_cmd(name):
    cfg = BIN_MATRIX[name]
    if cfg['type'] == 'go':
        if 'distro' in cfg.keys():
            for goos, archs in cfg['distro'].iteritems():
                for goarch in archs:
                    _go_build(name, goos, goarch)
        else:
            _go_build(name, GOHOSTOS, GOHOSTARCH)


def default():
    die(call('goimports -w pkg cmd plugins'))
    call('go fmt ./pkg/... ./cmd/... ./plugins/...')


if __name__ == "__main__":
    if len(sys.argv) > 1:
        globals()[sys.argv[1]](*sys.argv[2:])
    else:
        default()
