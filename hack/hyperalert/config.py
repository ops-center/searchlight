#!/usr/bin/env python

import subprocess
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
BACKEND_ROOT = expandvars('$GOPATH') + '/src/github.com/appscode/searchlight'
GOHOSTOS = GOENV["GOHOSTOS"]
GOHOSTARCH = GOENV["GOHOSTARCH"]
GOC = 'go'  # if ENV in ['dev'] else 'godep go'


# type must be specified
BIN_MATRIX = {
    'hyperalert': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_component_status': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_influx_query': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_json_path': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_node_count': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_node_status': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
    'check_pod_exists': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
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
            'linux': [
                'amd64'
            ]
        }
    },
    'check_volume': {
        'type': 'go',
        'go_version': True,
        'distro': {
            'linux': [
                'amd64'
            ]
        }
    },
}
