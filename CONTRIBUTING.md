# Contribution Guidelines for Aether

## Welcome!

We are happy to have contributors to our efforts to make IT sustainable. To get
started please check out any documentation including our [methodologies][1]

## Code Of Conduct

We adhere to the [contributors covenant][4]

the terms of this covenant will be subject to appropriate consequences, including but not limited to warnings, temporary bans, or permanent expulsion from the project community.

Our community values inclusivity, respect, and collaboration. We strive to create a welcoming environment for everyone, regardless of race, ethnicity, gender identity, sexual orientation, religion, disability, age, nationality, or any other personal characteristic.

We encourage constructive criticism, feedback, and discussions aimed at improving the project. However, we do not tolerate harassment, discrimination, personal attacks, or any form of disrespectful behavior.

We ask all contributors to:

1. Treat others with kindness, empathy, and respect.
2. Listen actively and considerately to diverse perspectives.
3. Avoid offensive language, derogatory remarks, or inappropriate imagery.
4. Respect differing opinions and gracefully resolve disagreements.
5. Use clear and inclusive language in all communications.
6. Refrain from engaging in any form of harassment, including but not limited to verbal, physical, sexual, or online harassment.
7. Be mindful of the impact of your words and actions on others.
8. Prioritize the well-being and safety of all community members.

By participating in this project, you agree to abide by these guidelines and uphold the principles of our community. Together, we can create a positive and inclusive environment where everyone feels valued and empowered to contribute.

### Development

To get direction on how to develop on Aether please see the [Development
Guide][5]

## How to Submit Changes

This repo adheres to conventional commits. Commits use the following structure

```log

<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Where types are:

- `fix`: designates a fix to the codebase
- `feat`: designates a new feature
- `chore`: house keeping tasks
- `docs`: documentation changes
- `ci`: changes to ci

Please make sure your commits adhere to this standard as release notes are
generated from this.

We dont accept merge commits, so please use rebasing.

A not on Pull Requests, please make sure that you changes are as small as
possible and within the scope of what you are trying to change. We will not
accept changes that are outside of the scope of the Pull Request.

for more information on conventional commits please see the [docs][2]


## Requesting an Enhancement

Please make sure all feature requests are logged as issues with a clear
explanation of the feature. 


## Golang Style guides

As much as possible we adhere to [googles golang style guide][3], please make sure
you are familiar with this guide before submitting code


[1]: ./methodologies.md
[2]: https://www.conventionalcommits.org/en/v1.0.0/
[3]: https://google.github.io/styleguide/go/
[4]: https://www.contributor-covenant.org/version/2/1/code_of_conduct/
[5]: ./DEVELOPMENT.md
