I want to create an app template that I can reuse over and over again for agentic development.
One example I want to use as a reference is located at ~/Personal/mark/*. The things that I really like from this setup is:
- Backend and frontend in one repository.
- Backend coding and testing standards:
  - Golang + Posgres as the foundational technology.
  - Testing standards where database access is included in the tests with pgflock and State creation and teardown.
  - Domain split of service, accessor, pservice, clear responsibilities.
  - How each domains are coded, init.sql files for each domain.
  - .env.example as the environment that is read when running the backend.
- Frontend coding and testing standards:
  - Using latest version of svelte.
  - Testing standards where we have unit tests and playtest, screenshot, manual testing scripts.
  - npm protected from supply chain attacks by minimum age when installing.
- Frontend hosting the PRD, storybook.
  - PRD and storybook only visible in development.
  - Allows writing PRD, storybook, decisions as we continue development.
  - Ensures context do not get lost between multiple agent sessions.
- Scripts directory:
  - Contains scripts to run dev backend, dev frontend, logs output at logs directory.
  - Plenty of scripts like health-checks, manual playwright test script, reset local db, init db, etc.
- Release directory:
  - Keeps track of migration scripts, release notes, changelog, etc. for each release.
  - Contains infrastructure code for deployment, nginx config, docker files, etc.
  - Point any agent to this repository and they can help release, iterate infrastructure,
    troubleshoot production issues, all context on production and staging environment is here.
- Repository .dev and .dev.example:
  - Allows multiple copies of the repository to run their own local backend and web instances.
  - We can clone the repo multiple times and have multiple agents running their tests and playwright tests,
    resetting local db etc. without impacting other agent sessions. Allows parallel development with tests
    and browser manual checks in one machine.
- AGENTS.md file that contains context of repository:
  - Coding standards, code architecture, do's and don'ts.
  - Every agent session opened in the workspace will not drift away from quality standards and knows
    the architecture and where to find things, where to write, how to run test, etc.

The goal is we can quickly spin up an app repository and do agentic engineering that results in
the same code architecture quality, code quality, testing, documentation, deployment standards that we have
in mark, so we don't need to reinvent the wheel every time.

Go through ~/Personal/mark/* and figure out everything that we can extract out into this template.
Other than what I mentioned above, also look for any other practices, patterns, architecture decisions,
coding styles, testing styles, documentation styles, deployment styles that we can extract out into this template.