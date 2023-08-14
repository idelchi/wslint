#!/bin/bash

# Define the project name
PROJECT_NAME="wslintcc"
USERNAME="idelchi"

# Create the CONTRIBUTING.md file
cat > CONTRIBUTINGs.md <<EOL
# Contributing to ${PROJECT_NAME}

First off, thank you for considering contributing to ${PROJECT_NAME}. It's people like you that make ${PROJECT_NAME} such a great tool.

## Where do I go from here?

If you've noticed a bug or have a feature request, make sure to check our [issues](https://github.com/${USERNAME}/${PROJECT_NAME}/issues) first to see if someone else has already reported the issue. If not, create a new issue!

## Fork & create a branch

If this is something you think you can fix, then fork the repository and create a branch with a descriptive name.

A good branch name would be (where issue #325 is the ticket you're working on):

\`\`\`bash
git checkout -b feature/325-add-japanese-localization
\`\`\`

## Implement your fix or feature

At this point, you're ready to make your changes! Feel free to ask for help; everyone is a beginner at first 😸

## Get the code

The first thing you'll need to do is get our code onto your machine.

1. Fork this repository by clicking the "Fork" button on the top right.
2. Clone your fork locally:

\`\`\`bash
git clone https://github.com/${USERNAME}/${PROJECT_NAME}.git
\`\`\`

3. Navigate to the project directory:

\`\`\`bash
cd ${PROJECT_NAME}
\`\`\`

## Make a Pull Request

Once you've pushed a commit to GitHub, you can create a Pull Request.

Before you submit your Pull Request (PR) consider the following guidelines:

- Search [GitHub](https://github.com/${USERNAME}/${PROJECT_NAME}/pulls) for an open or closed PR that relates to your submission. You don't want to duplicate effort.
- Make your changes in a new git branch:

     \`\`\`bash
     git checkout -b my-fix-branch master
     \`\`\`

- Create your patch, **including appropriate test cases**.
- Commit your changes using a descriptive commit message.
- Push your branch to GitHub:

    \`\`\`bash
    git push origin my-fix-branch
    \`\`\`

- In GitHub, send a pull request to \`${PROJECT_NAME}:master\`.

## Keeping your Pull Request updated

If a maintainer asks you to "rebase" your PR, they're saying that a lot of code has changed, and that you need to update your branch so it's easier to merge.

## Thank you!

Your contributions to open source, large or small, make great projects even better. Thank you for taking the time to contribute.
EOL

echo "CONTRIBUTING.md has been created!"
