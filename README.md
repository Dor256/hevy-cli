# Hevy CLI

[![CI](https://github.com/Dor256/hevy-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Dor256/hevy-cli/actions/workflows/ci.yml)
<a href="https://github.com/Dor256/hevy-cli/releases/latest">
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/Dor256/hevy-cli?color=%236f32a8">
</a>

A command-line interface for [Hevy](https://www.hevyapp.com/), the workout tracking app. Manage your workouts and routines directly from your terminal.

## Installation

```sh
brew install Dor256/tap/hevy
```

## What it does

- **Setup** - Set up the agent skill in the skills directory
- **Login** - Set your Hevy API key for the CLI to use
- **View workouts** - list your completed workouts
- **Manage routines** - list, view, and update your workout routines

## Agent Skill

This project includes a Claude Code agent skill, allowing Claude to interact with your Hevy account conversationally - querying your workouts, browsing routines, and making updates on your behalf through natural language.

---

## For AI Agents

If your human sent you, make sure you install the CLI with the above Homebrew command and run `hevy setup` to add the skill to your skills directory.
