Genesis CLI
======

## Installation

### From source

On all operating systems, this CLI can be installed with the command

```
make install
```

### Without go tools

#### Linux and Mac

```
curl -sSf https://assets.whiteblock.io/cli/install.sh | sh
```

#### Windows
Download [here](https://assets.whiteblock.io/cli/master/bin/windows/amd64/genesis.exe)

## Running a test
To run a test as a user in the org named "whiteblock", you would use the command:

```
genesis run <path to your yaml spec> <your-organization-name>
```

If it is your first time using the CLI, it will prompt you in your browser to complete the authentication process. You only need to do this once.


The CLI will also remember the last organization you used, so after running that command, you would only need to use

```
genesis run <path to your yaml spec>
```

to run a subsequent test in the org "whiteblock".

### DNS
To easily access your running test, you can choose for the instance hosts of the test to be exposed using DNS.


This is done by adding the `--dns` flag to your run command, like so:

```
genesis run --dns ...
```

The result will be the useable DNS names printed out for you on a per test basis. Here is an example of what it might look like:
```
success
Definition: fake-definition-id

Test: your-test
	ratbog-0.biomes.whiteblock.io
	ratbog-1.biomes.whiteblock.io
Test: your-other-test
	wingtree-0.biomes.whiteblock.io
	wingtree-1.biomes.whiteblock.io

```
