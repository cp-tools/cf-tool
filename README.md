<h1 align="center">
    <img src="assets/logo.png" alt="logo" width="20%" height="20%">
    <br/>
    Codeforces Tool
</h1>
<h4 align="center">
    An amazingly quick, lightweight, yet powerful tool for <a href="https://codeforces.com">CodeForces</a>
    <br />
    Don't forget to :star: the project if you liked it!
    <br /><br />
    <a href="https://travis-ci.com/github/infixint943/cf"><img src="https://img.shields.io/travis/com/infixint943/cf?style=for-the-badge"></a>
    <a href="https://github.com/infixint943/cf/commits/master"><img src="https://img.shields.io/github/last-commit/infixint943/cf?style=for-the-badge"></a>
    <a href="https://github.com/infixint943/cf/releases"><img src="https://img.shields.io/github/v/release/infixint943/cf?style=for-the-badge"></a>
    <a href="https://github.com/infixint943/cf/issues"><img src="https://img.shields.io/github/issues/infixint943/cf?style=for-the-badge"></a>
    <a href=""><img src="https://img.shields.io/github/go-mod/go-version/infixint943/cf?style=for-the-badge"></a>
    <a href="https://github.com/infixint943/cf/blob/master/LICENSE"><img src="https://img.shields.io/github/license/infixint943/cf?style=for-the-badge"></a>
</h4>



![](assets/demo.gif)



# Key Features

- Fetch test cases of entire contests or individual problems, to well structured directories.
- Supports official contests, gym contests as well as groups.
- Compile and run source code (locally) against test cases.
- Set custom timeout to prevent system hang.
- Manual testing of interactive problems.
- Validation of output using custom checkers.
- Submit solutions directly and view (dynamic) status of submission.
- Pull submission(s) of any particular user.
- Ability to configure mirror domain, proxy protocols.
- Beautifully crafted, colourful CLI.

# Installation

You may download the latest, compiled, binary files from [here](https://github.com/infixint943/cf/releases).
Place the executable in system **PATH** to invoke the tool from any directory.

Alternatively, you can also compile the tool from source.

```bash
git clone https://github.com/infixint943/cf.git
cd cf/
go build -ldflags "-s -w"
```

# Quick Start

**Note:** For detailed documentation, please head to the [wiki](https://github.com/infixint943/cf/wiki) page.

> Let's simulate participating in contest `4`. This tutorial assumes you have already configured your login details and added at least one template, through `cf config` 

- To start competing, run `cf fetch 4`. If the contest hasn't started, a countdown will be launched.
- Once the contest begins, test cases of every problem in contest `4` are fetched to the current folder like in the folder structure below. Also, the [problems](https://codeforces.com/contest/4/problems) page will automatically be opened in your default browser.
- Navigate to the folder containing problem `d`. This folder will contain all sample inputs (`*.in`), sample outputs (`*.out`) and source code files.
- Generate a template source file (if not generated automatically) with command `cf gen`. The source file will be saved with file name `d`.
- Edit the generated source file with any tool/IDE of your choice.
- When you wish to test your solution against the sample input files, run command `cf test`. Your source file will be auto determined and will be run against all test cases. The verdict of each test case will be displayed.
- Once you want to submit your source file, run `cf submit`. Your source file will be (auto determined and) submitted to the respective problem (again, determined from the folder path) and the verdict will be dynamically updated in the terminal.
- To view the contest dashboard (solve count of every problem in the contest) run `cf watch 4`. The specified details with the current status will be displayed in the terminal.

# FAQ

### How do I install and use this tool?
Place the compiled binary file in system `PATH`. Then you can run the tool from the terminal with command `cf`

### How is this different from xalanq/cf-tool?

This tool actually draws most of its inspiration from [xalanq/cf-tool](https://github.com/xalanq/cf-tool). However, this tool has several changes, of which the most notable ones are:

- An overhauled, more polished UI
- Unification of redundant commands (`pull` and `clone`, `parse` and `race`)
- New custom checkers for testing - case insensitivity, ignore difference in decimal output after certain EPS, custom time limit.
- Improved code logic, leading to faster fetching of sample tests.
- [WIP] API functionality to return plain text output (without escape sequence).
- A much better line-by-line diff comparision (as compared to the traditional diff tool, which isn't very helpful to compare multi line outputs) 

Apart from these changes, the code has been restructured entirely, resulting in a much shorter and clean code base (helpful for contributors, as everything is well commented in the project).  

### Where are the configuration files stored?

The configuration files are stored in the folder `cf`, created in the `$XDG_CONFIG_DIR` directory. The path to the directory containing these configuration files are:

| Operating System | Folder Path                            |
| ---------------- | -------------------------------------- |
| Darwin           | `$HOME/Library/Application Support/cf` |
| Windows          | `%AppData%/cf`                         |
| Linux            | `$HOME/.config/cf`                     |

If `$XDG_CONFIG_DIR` path can't be determined, the configurations will be saved in folder `$HOME/cf` folder.

### How do I enable tab completion in the terminal?

You can download the completion scripts (bash and zsh) provided along with every release. [Configure](https://stackoverflow.com/questions/45115260/where-to-put-bash-completion-scripts) it suitably according to your OS. Windows user might need to follow additional steps to first enable tab completion in the terminal.

Alternatively, you can try following instructions in [docopt_completion](https://github.com/Infinidat/infi.docopt_completion) to automatically add the completion scripts to the required system settings.

**Note:** If the man page changes upon subsequent release, you may have to reconfigure the tab completion.
