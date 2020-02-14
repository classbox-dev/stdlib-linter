from setuptools import setup

setup(
    name='stdlib-linters',
    packages=['linters'],
    entry_points={'console_scripts': [
        'pygofmt = linters.gofmt:main',
        'pyuntracked = linters.untracked:main',
    ]},
)
