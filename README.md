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
    <a href=""><img src="https://img.shields.io/travis/com/infixint943/cf?style=for-the-badge"></a>
    <a href=""><img src="https://img.shields.io/github/last-commit/infixint943/cf?style=for-the-badge"></a>
    <a href=""><img src="https://img.shields.io/github/v/release/infixint943/cf?style=for-the-badge"></a>
    <a href=""><img src="https://img.shields.io/github/issues/infixint943/cf?style=for-the-badge"></a>
    <a href=""><img src="https://img.shields.io/github/go-mod/go-version/infixint943/cf?style=for-the-badge"></a>
    <a href=""><img src="https://img.shields.io/github/license/infixint943/cf?style=for-the-badge"></a>
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

**Note:** `bash` and `zsh` auto completion scripts are provided in [this tarball]().

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
