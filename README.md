# alfred-project-manager

> **v1.6.0**
> 
> - 支持冒号分隔的多项目目录（`~/Code:~/Projects`）
> - 智能项目检测：项目标记文件 + vendor 目录 + 源码文件，自动过滤非项目目录
> - 修饰键重构：回车→Finder / shift→ClaudeCode / cmd→VSCode / opt→Kitty
> - 默认参数可选、无需空格触发、立即响应
> 
> Allows you to quickly open projects from Alfred.

![usage example](/image.png)

## Installation

- Grab the latest release [here](https://github.com/bjrnt/alfred-project-manager/releases/) and install the workflow file.

## Usage

Open Alfred and type `pm` to access the project manager and try typing a query. You can also configure a hotkey for it by opening the workflow in Alfred's workflow options panel.

### Modifiers

- `none` (default): reveal in Finder
- `shift`: open in Kitty with ClaudeCode (`claude /path/to/project`)
- `cmd`: open in VSCode
- `opt`: open in Kitty
- `ctrl`: open the project's repo in your browser

You can change the modifier and application combinations in Alfred's workflow settings window.

### Features

- **Multiple Project Directories** — Configure multiple project collection paths separated by `:` (e.g. `~/Code:~/Projects`). Each path's immediate children are scanned as project candidates. Avoid broad directories like `~/Desktop`.
- **Detect Code Projects** — Enable intelligent project detection. Uses a three-layer check: project marker files (package.json, go.mod, .git, etc.), vendor/build directories (node_modules, dist, etc.), and source code files. Directories without any code artifacts are filtered out.
- **Maximum Project Depth** — Controls how many levels deep to scan for projects. Defaults to 999 (practically unlimited). Set to 0 to scan only immediate children of your configured directories.

## Maintenance

### Issues and Feature Requests

Feel free to open an issue for the project if you have encountered a problem or have a feature request for the workflow.

### Building

The project can be built, linked to Alfred, and released using [jason0x34/go-alfred](https://github.com/jason0x43/go-alfred). Commands for this can be found in [./.vscode/tasks.json](./.vscode/tasks.json).

## Changelog

### v1.6.0

- Support multiple project directories (colon-separated: `~/Code:~/Projects`)
- Smart project detection via go-enry (project markers, vendor dirs, source code files)
- `Detect Code Projects` checkbox replaces `Require Git Repo`
- Modifier keys: none→Finder, shift→ClaudeCode, cmd→VSCode, opt→Kitty
- Default: withspace off, argument optional, immediate queue, depth 999
- Subtitles display full project path with `~` for $HOME
