# -*- coding: utf-8 -*-
from __future__ import absolute_import
from __future__ import unicode_literals

import pytest

import dockerfile


def test_command_module():
    """namedtuple defaults to the parent python frame. we fix this"""
    assert dockerfile.Command.__module__ == dockerfile.__name__


def test_all_cmds():
    assert dockerfile.all_cmds()[:3] == ('add', 'arg', 'cmd')


def test_parse_file_ioerror():
    with pytest.raises(dockerfile.GoIOError) as excinfo:
        dockerfile.parse_file('Dockerfile.dne')
    assert 'Dockerfile.dne' in excinfo.value.args[0]


def test_parse_string_parse_error():
    with pytest.raises(dockerfile.GoParseError):
        dockerfile.parse_string(
            'FROM ubuntu:xenial\n'
            'CMD ["echo", 1]\n'
        )


def test_parse_string_success():
    ret = dockerfile.parse_string(
        'FROM ubuntu:xenial\n'
        'RUN echo hi > /etc/hi.conf\n'
        'CMD ["echo"]\n'
        'HEALTHCHECK --retries=5 CMD echo hi\n'
        'ONBUILD ADD foo bar\n'
        'ONBUILD RUN ["cat", "bar"]\n'
    )
    assert ret == (
        dockerfile.Command(
            cmd='from', sub_cmd=None, json=False, flags=(),
            value=('ubuntu:xenial',),
            start_line=1, original='FROM ubuntu:xenial',
        ),
        dockerfile.Command(
            cmd='run', sub_cmd=None, json=False, flags=(),
            value=('echo hi > /etc/hi.conf',),
            start_line=2, original='RUN echo hi > /etc/hi.conf',
        ),
        dockerfile.Command(
            cmd='cmd', sub_cmd=None, json=True, flags=(), value=('echo',),
            start_line=3, original='CMD ["echo"]',
        ),
        dockerfile.Command(
            cmd='healthcheck', sub_cmd=None, json=False,
            flags=('--retries=5',), value=('CMD', 'echo hi'),
            start_line=4, original='HEALTHCHECK --retries=5 CMD echo hi',
        ),
        dockerfile.Command(
            cmd='onbuild', sub_cmd='add', json=False, flags=(),
            value=('foo', 'bar'),
            start_line=5, original='ONBUILD ADD foo bar',
        ),
        dockerfile.Command(
            cmd='onbuild', sub_cmd='run', json=True, flags=(),
            value=('cat', 'bar'),
            start_line=6, original='ONBUILD RUN ["cat", "bar"]',
        ),
    )


def test_parse_string_text():
    ret = dockerfile.parse_string(
        'FROM ubuntu:xenial\n'
        'CMD ["echo", "☃"]\n'
    )
    assert ret == (
        dockerfile.Command(
            cmd='from', sub_cmd=None, json=False, value=('ubuntu:xenial',),
            start_line=1, original='FROM ubuntu:xenial', flags=(),
        ),
        dockerfile.Command(
            cmd='cmd', sub_cmd=None, json=True, value=('echo', '☃'),
            start_line=2, original='CMD ["echo", "☃"]', flags=(),
        ),
    )


def test_parse_file_success():
    ret = dockerfile.parse_file('testfiles/Dockerfile.ok')
    assert ret == (
        dockerfile.Command(
            cmd='from', sub_cmd=None, json=False, flags=(),
            value=('ubuntu:xenial',),
            start_line=1, original='FROM ubuntu:xenial',
        ),
        dockerfile.Command(
            cmd='cmd', sub_cmd=None, json=True, flags=(), value=('echo', 'hi'),
            start_line=2, original='CMD ["echo", "hi"]',
        ),
    )
