from setuptools import setup

setup(
    name='stdlib-hooks',
    packages=['hooks'],
    entry_points={'console_scripts': [
        'pygofmt = hooks.gofmt:main',
        'pyuntracked = hooks.untracked:main',
    ]},
)
