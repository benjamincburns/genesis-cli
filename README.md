======

## Installation

### Binaries

#### Linux and Mac

```
curl -sSf https://assets.whiteblock.io/cli/install.sh | sh
```

#### Windows
Download [here](https://assets.whiteblock.io/cli/bin/windows/amd64/genesis.exe)

### From source

On all operating systems, this CLI can be installed with the command

```
make install
```

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

to run a subsequent test in the org.


### Remote Auth
You can get your credentials with the command `genesis auth print-access-token`
And then you can authenticate anywhere by setting the environment variable GENESIS_CREDENTIALS to the value of the previous command.
