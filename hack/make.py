#!/usr/bin/env python

import datetime
import io
import json
import os
import os.path
import socket
import subprocess
import sys
from collections import OrderedDict
from os.path import expandvars


def _goenv():
    env = {}
    for line in subprocess.check_output(['go', 'env']).split('\n'):
        line = line.strip()
        if len(line) == 0:
            continue
        k, v = line.split('=', 1)
        v = v.strip('"')
        if len(v) > 0:
            env[k] = v
    return env


GOENV = _goenv()
GOHOSTOS = GOENV["GOHOSTOS"]
GOHOSTARCH = GOENV["GOHOSTARCH"]
GOC = 'go'  # if ENV in ['dev'] else 'godep go'


def metadata(cwd, goos='', goarch=''):
    # Nearest tag
    # try:
    #     tag = subprocess.check_output('git describe --abbrev=0 2>/dev/null', shell=True, cwd=cwd).strip()
    # except subprocess.CalledProcessError:
    #     pass

    metadata = {
        'commit_hash': subprocess.check_output('git rev-parse --verify HEAD', shell=True, cwd=cwd).strip(),
        'git_branch': subprocess.check_output('git rev-parse --abbrev-ref HEAD', shell=True, cwd=cwd).strip(),
        # http://stackoverflow.com/a/1404862/3476121
        'git_tag': subprocess.check_output('git describe --exact-match --abbrev=0 2>/dev/null || echo ""', shell=True,
                                           cwd=cwd).strip(),
        'commit_timestamp': datetime.datetime.utcfromtimestamp(
            int(subprocess.check_output('git show -s --format=%ct', shell=True, cwd=cwd).strip())).isoformat(),
        'build_timestamp': datetime.datetime.utcnow().isoformat(),
        'build_host': socket.gethostname(),
        'build_host_os': GOENV["GOHOSTOS"],
        'build_host_arch': GOENV["GOHOSTARCH"]
    }
    if metadata['git_tag']:
        metadata['version'] = metadata['git_tag']
        metadata['version_strategy'] = 'tag'
    elif not metadata['git_branch'] in ['master', 'HEAD']:
        metadata['version'] = metadata['git_branch']
        metadata['version_strategy'] = 'branch'
    else:
        commit_ts = subprocess.check_output('git show -s --format=%ct', shell=True, cwd=cwd).strip()
        metadata['version'] = datetime.datetime.utcfromtimestamp(int(commit_ts)).strftime('%Y%m%d')
        # metadata['version'] = subprocess.check_output('TZ=UTC gdate -d @$(git show -s --format=%ct) +"%Y%m%d"', shell=True, cwd=cwd).strip()
        nearest_tag = subprocess.check_output('git describe --abbrev=0 2>/dev/null || echo "0.0"', shell=True,
                                              cwd=cwd).strip()
        if nearest_tag:
            metadata['version'] = nearest_tag + '.' + metadata['version']
        metadata['version_strategy'] = 'timestamp'
    if goos:
        metadata['os'] = goos
    if goarch:
        metadata['arch'] = goarch
    return metadata


# Debian package
# https://gist.github.com/rcrowley/3728417
REPO_ROOT = expandvars('$GOPATH') + '/src/github.com/appscode/searchlight'
ENV = os.getenv('APPSCODE_ENV', 'dev').lower()
BUILD_METADATA = metadata(REPO_ROOT)
BIN_MATRIX = {
    'hyperalert': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_component_status': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_influx_query': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_json_path': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_node_count': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_node_status': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_pod_exists': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_pod_status': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_prometheus_metric': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
    'check_volume': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': ['amd64']
        }
    },
}
BUCKET_MATRIX = {
    'prod': 'gs://appscode-cdn',
    'dev': 'gs://appscode-dev'
}


def read_file(name):
    with open(name, 'r') as f:
        return f.read()
    return ''


def write_file(name, content):
    dir = os.path.dirname(name)
    if not os.path.exists(dir):
        os.makedirs(dir)
    with open(name, 'w') as f:
        return f.write(content)


def append_file(name, content):
    with open(name, 'a') as f:
        return f.write(content)


def write_checksum(folder, file):
    cmd = "openssl md5 {0} | sed 's/^.* //' > {0}.md5".format(file)
    subprocess.call(cmd, shell=True, cwd=folder)
    cmd = "openssl sha1 {0} | sed 's/^.* //' > {0}.sha1".format(file)
    subprocess.call(cmd, shell=True, cwd=folder)


# TODO: use unicode encoding
def read_json(name):
    try:
        with open(name, 'r') as f:
            return json.load(f, object_pairs_hook=OrderedDict)
    except IOError:
        return {}


def write_json(obj, name):
    with io.open(name, 'w') as f:
        data = json.dumps(obj, indent=2, separators=(',', ': '), ensure_ascii=False)
        f.write(data)


def call(cmd, stdin=None, cwd=REPO_ROOT):
    print(cmd)
    return subprocess.call([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def die(status):
    if status:
        sys.exit(status)


def check_output(cmd, stdin=None, cwd=REPO_ROOT):
    print(cmd)
    return subprocess.check_output([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def version():
    # json.dump(BUILD_METADATA, sys.stdout, sort_keys=True, indent=2)
    for k in sorted(BUILD_METADATA):
        print k + '=' + BUILD_METADATA[k]


def fmt():
    die(call('goimports -w pkg cmd plugins'))
    call('go fmt ./pkg/... ./cmd/... ./plugins/...')


def vet():
    call('go vet ./pkg/... ./cmd/... ./plugins/...')


def gen():
    # gen_assets()
    # gen_extpoints()
    return


def build_cmd(name):
    cfg = BIN_MATRIX[name]
    if cfg['type'] == 'go':
        if 'distro' in cfg.keys():
            for goos, archs in cfg['distro'].iteritems():
                for goarch in archs:
                    _go_build(name, goos, goarch)
        else:
            _go_build(name, GOHOSTOS, GOHOSTARCH)


def _to_upper_camel(lower_snake):
    components = lower_snake.split('_')
    # We capitalize the first letter of each component
    # with the 'title' method and join them together.
    return ''.join(x.title() for x in components[:])


# ref: https://golang.org/cmd/go/
def _go_build(name, goos, goarch):
    linker_opts = []
    if BIN_MATRIX[name].get('go_version', False):
        for k, v in metadata(REPO_ROOT, goos, goarch).iteritems():
            linker_opts.append('-X')
            linker_opts.append('main.' + _to_upper_camel(k) + '=' + v)

    cgo_env = 'CGO_ENABLED=0'
    cgo = ''
    if BIN_MATRIX[name].get('use_cgo', False):
        cgo_env = "CGO_ENABLED=1"
        cgo = "-a -installsuffix cgo"
        linker_opts.append('-linkmode external -extldflags -static -w')

    ldflags = ''
    if linker_opts:
        ldflags = "-ldflags '{}'".format(' '.join(linker_opts))

    bindir = 'dist/{name}'.format(name=name)
    if not os.path.isdir(bindir):
        os.makedirs(bindir)
    cmd = "GOOS={goos} GOARCH={goarch} {cgo_env} {goc} build -o {bindir}/{name}-{goos}-{goarch}{ext} {cgo} {ldflags} cmd/{name}/*.go".format(
        name=name,
        goc=GOC,
        goos=goos,
        goarch=goarch,
        bindir=bindir,
        cgo_env=cgo_env,
        cgo=cgo,
        ldflags=ldflags,
        ext='.exe' if goos == 'windows' else ''
    )
    die(call(cmd))
    print '\n'


def build_cmds():
    gen()
    fmt()
    for name in BIN_MATRIX.keys():
        build_cmd(name)


def build_all():
    build_cmds()


def push_all():
    dist = REPO_ROOT + '/dist'
    for name in os.listdir(dist):
        d = dist + '/' + name
        if os.path.isdir(d):
            push(d)


def push_cmd(name):
    bindir = REPO_ROOT + '/dist/' + name
    push(bindir)


def push(bindir):
    call('rm -f *.md5', cwd=bindir)
    call('rm -f *.sha1', cwd=bindir)
    for f in os.listdir(bindir):
        if os.path.isfile(bindir + '/' + f):
            _upload_to_cloud(bindir, f)


def _upload_to_cloud(folder, f):
    write_checksum(folder, f)
    name = os.path.basename(folder)
    if name not in BIN_MATRIX.keys():
        return
    if ENV == 'prod' and not BIN_MATRIX[name].get('release', False):
        return

    bucket = BUCKET_MATRIX.get(ENV, BUCKET_MATRIX['dev'])
    dst = "{bucket}/binaries/{name}/{version}/{file}{ext}".format(
        bucket=bucket,
        name=name,
        version=BUILD_METADATA['version'],
        file=f,
        ext='.exe' if '-windows-' in f else ''
    )
    if bucket.startswith('gs://'):
        _upload_to_gcs(folder, f, dst, BIN_MATRIX[name].get('release', False))


def _upload_to_gcs(folder, src, dst, public):
    call("gsutil cp {0} {1}".format(src, dst), cwd=folder)
    call("gsutil cp {0}.md5 {1}.md5".format(src, dst), cwd=folder)
    call("gsutil cp {0}.sha1 {1}.sha1".format(src, dst), cwd=folder)
    if public:
        call("gsutil acl ch -u AllUsers:R {0}".format(dst), cwd=folder)
        call("gsutil acl ch -u AllUsers:R {0}.md5".format(dst), cwd=folder)
        call("gsutil acl ch -u AllUsers:R {0}.sha1".format(dst), cwd=folder)


def update_registry():
    dist = REPO_ROOT + '/dist'
    bucket = BUCKET_MATRIX.get(ENV, BUCKET_MATRIX['dev'])
    lf = dist + '/latest.txt'
    write_file(lf, BUILD_METADATA['version'])
    for name in os.listdir(dist):
        if name not in BIN_MATRIX.keys():
            return
        call("gsutil cp {2} {0}/binaries/{1}/latest.txt".format(bucket, name, lf))
        if BIN_MATRIX[name].get('release', False):
            call('gsutil acl ch -u AllUsers:R -r {0}/binaries/{1}/latest.txt'.format(bucket, name))


def install():
    die(call('GO15VENDOREXPERIMENT=1 ' + GOC + ' install ./pkg/... ./cmd/... ./plugins/...'))


def default():
    gen()
    fmt()
    die(call('GO15VENDOREXPERIMENT=1 ' + GOC + ' install ./pkg/... ./cmd/... ./plugins/...'))


if __name__ == "__main__":
    if len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        globals()[sys.argv[1]](*sys.argv[2:])
    else:
        default()
