import os

from setuptools import Extension
from setuptools import setup


if not os.path.exists('vendor/github.com/docker/docker/builder'):
    print('docker checkout is missing!')
    print('Run `git submodule update --init`')
    exit(1)


setup(
    name='dockerfile',
    description=(
        'Parse a dockerfile into a high-level representation using the '
        'official go parser.'
    ),
    url='https://github.com/asottile/dockerfile',
    version='2.0.0',
    author='Anthony Sottile',
    author_email='asottile@umich.edu',
    classifiers=[
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 2',
        'Programming Language :: Python :: 2.7',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: Implementation :: CPython',
        'Programming Language :: Python :: Implementation :: PyPy',
    ],
    ext_modules=[
        Extension('dockerfile', ['pylib/main.go']),
    ],
    build_golang={'root': 'github.com/asottile/dockerfile'},
    setup_requires=['setuptools-golang>=0.2.0'],
)
