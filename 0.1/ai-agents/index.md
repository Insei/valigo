# For AI agents

AI agents often write the code that uses a library, and a model trained before the typed configurator doesn't know its API: it writes `b.String(&u.Name)` of the old API instead of `b.Field(&u.Name).As(str.Of)`, or `Max(int64(5))` with an `any` argument. Give the agent the current documentation.

## The skill

[`skills/valigo/SKILL.md`](https://github.com/Insei/valigo/blob/main/skills/valigo/SKILL.md) is a short guide for agents in the [Agent Skills](https://agentskills.io) format: the rules, complete programs, the configurators with their rules and the common mistakes. It is versioned with the code, and CI builds and runs its programs against the code of the repository, so it doesn't describe the API that doesn't exist.

Claude Code finds the skill in the `.claude/skills` directory of the project or of the user:

```sh
mkdir -p .claude/skills/valigo
curl -fsSL https://raw.githubusercontent.com/Insei/valigo/main/skills/valigo/SKILL.md \
  -o .claude/skills/valigo/SKILL.md
```

Replace `main` with the release tag to get the skill of the version in your `go.mod`. Other agents read it from any path: mention it in `AGENTS.md` or in the instructions of the agent.

The skill describes the library, not your project. Keep the rules of your project — where the validators are configured, how the errors are returned to the clients — in your own instructions next to it.

## llms.txt

The documentation is published in the [llms.txt](https://llmstxt.org) format for every version:

- [`llms.txt`](https://insei.github.io/valigo/latest/llms.txt) — the index of the pages with descriptions;
- [`llms-full.txt`](https://insei.github.io/valigo/latest/llms-full.txt) — the whole documentation as one Markdown file, to put into the context of the agent.

Replace `latest` with the version, or with `dev` for the documentation of the `main` branch.
